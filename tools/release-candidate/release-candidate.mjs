#!/usr/bin/env node

import { createHash } from 'node:crypto';
import { execFileSync } from 'node:child_process';
import { readFileSync, writeFileSync } from 'node:fs';

export const metadataPath = '.nx/release-candidate.json';

export function createReleaseMetadata({ source = 'HEAD^', candidate = 'HEAD' } = {}) {
  const planPaths = git(['ls-tree', '-r', '--name-only', source, '--', '.nx/version-plans'])
    .split('\n')
    .filter(path => path.endsWith('.md'));
  if (planPaths.length === 0) fail('Reviewed source has no version plans.');

  const requested = new Map();
  const plans = planPaths.map(path => {
    const content = git(['show', `${source}:${path}`], false);
    for (const [project, bump] of parseVersionPlan(content, path)) {
      const previous = requested.get(project);
      requested.set(project, strongestBump(previous, bump));
    }
    return { path, sha256: sha256(content) };
  });

  const releases = [...requested.entries()].sort(([a], [b]) => a.localeCompare(b)).map(
    ([project, bump]) => releaseForProject({ project, bump, source, candidate }),
  );
  const plannedTags = new Set(releases.map(release => release.npmTag));
  const candidateTags = git(['tag', '--points-at', candidate])
    .split('\n')
    .filter(tag => /-contract-v\d+\.\d+\.\d+$/.test(tag));
  const unplannedTags = candidateTags.filter(tag => !plannedTags.has(tag));
  if (unplannedTags.length) {
    fail(`Generated candidate contains unplanned npm tags: ${unplannedTags.join(', ')}`);
  }
  for (const release of releases) {
    if (!candidateTags.includes(release.npmTag)) {
      fail(`Generated candidate is missing ${release.npmTag}.`);
    }
  }

  return {
    schemaVersion: 1,
    reviewedSource: resolveCommit(source),
    plans,
    releases,
  };
}

export function validateFinalCandidate({ checked = 'HEAD', metadata } = {}) {
  const checkedCommit = resolveCommit(checked);
  const value = metadata ?? readMetadata(checkedCommit);
  validateMetadataShape(value);
  const metadataCommits = git([
    'log',
    '--format=%H',
    `${value.reviewedSource}..${checkedCommit}`,
    '--',
    metadataPath,
  ]).split('\n').filter(Boolean);
  if (metadataCommits.length !== 1) {
    fail(`Checked main must contain exactly one generated release candidate, found ${metadataCommits.length}.`);
  }
  const [candidate] = metadataCommits;
  if (resolveCommit(`${candidate}^`) !== value.reviewedSource) {
    fail('Generated candidate parent differs from its reviewed source.');
  }
  ensureAncestor(value.reviewedSource, candidate, 'reviewed source is not the candidate parent');
  ensureAncestor(candidate, checkedCommit, 'generated candidate is not established on checked main');

  const parents = git(['rev-list', '--parents', '-n', '1', checkedCommit]).split(/\s+/);
  if (
    parents.length !== 3 ||
    parents[1] !== value.reviewedSource ||
    parents[2] !== candidate
  ) {
    fail('Checked main must be a two-parent merge whose second parent is the generated candidate.');
  }

  if (git(['rev-parse', `${candidate}^{tree}`]) !== git(['rev-parse', `${checkedCommit}^{tree}`])) {
    fail('Checked main tree differs from the exact generated candidate tree.');
  }

  const candidateMetadata = readMetadata(candidate);
  if (stableJson(candidateMetadata) !== stableJson(value)) {
    fail('Release candidate metadata changed after generation.');
  }
  validateGeneratedTree({ source: value.reviewedSource, candidate, metadata: value });

  const protectedPaths = [metadataPath];
  for (const release of value.releases) {
    protectedPaths.push(release.manifestPath, release.changelogPath);
  }
  for (const path of protectedPaths) {
    if (git(['show', `${candidate}:${path}`], false) !== git(['show', `${checkedCommit}:${path}`], false)) {
      fail(`${path} changed after the generated candidate.`);
    }
  }
  if (git(['ls-tree', '-r', '--name-only', checkedCommit, '--', '.nx/version-plans']).trim()) {
    fail('Checked main still contains version plans.');
  }

  return { ...value, generatedCandidate: candidate, checkedMain: checkedCommit };
}

export function validateCiRevision({ event, revision = 'HEAD', base } = {}) {
  const checked = resolveCommit(revision);
  if (event === 'pull_request') {
    if (!base) fail('Pull-request candidate validation requires its checked base.');
    const checkedBase = resolveCommit(base);
    if (sameFile(checkedBase, checked, metadataPath)) {
      validateOrdinaryRevision(checkedBase, checked);
      return { mode: 'ordinary' };
    }
    const metadata = readMetadata(checked);
    validateGeneratedTree({ source: metadata.reviewedSource, candidate: checked, metadata });
    if (metadata.reviewedSource !== checkedBase) {
      fail('Generated release pull request does not directly target its reviewed source.');
    }
    return { mode: 'generated-candidate', ...metadata, generatedCandidate: checked };
  }
  if (event === 'push') {
    const parent = resolveCommit(`${checked}^`);
    if (sameFile(parent, checked, metadataPath)) {
      validateOrdinaryRevision(parent, checked);
      return { mode: 'ordinary' };
    }
    return { mode: 'checked-main', ...validateFinalCandidate({ checked }) };
  }
  fail(`Unsupported CI event: ${event}.`);
}

export function classifyPublication({
  releases,
  checkedMain,
  remoteTags,
  npmIntegrities,
  localIntegrities,
}) {
  const publishProjects = [];
  const tagsToPush = [];
  for (const release of releases) {
    const remoteTag = remoteTags[release.npmTag];
    if (remoteTag && (!remoteTag.rawObject || !remoteTag.peeledCommit)) {
      fail(`Remote tag ${release.npmTag} must be an annotated tag with a peeled commit.`);
    }
    if (remoteTag?.peeledCommit && remoteTag.peeledCommit !== checkedMain) {
      fail(`Remote tag ${release.npmTag} resolves to ${remoteTag.peeledCommit}, not checked main ${checkedMain}.`);
    }
    if (!remoteTag) tagsToPush.push(release.npmTag);
    const packageVersion = `${release.package}@${release.version}`;
    const localIntegrity = localIntegrities[packageVersion];
    if (!/^sha(256|384|512)-\S+$/.test(localIntegrity ?? '')) {
      fail(`Invalid packed integrity for ${packageVersion}.`);
    }
    const integrity = npmIntegrities[packageVersion];
    if (integrity === undefined || integrity === null || integrity === '') {
      publishProjects.push(release.project);
    } else if (!/^sha(256|384|512)-\S+$/.test(integrity)) {
      fail(`Invalid npm integrity for ${packageVersion}.`);
    } else if (integrity !== localIntegrity) {
      fail(`Registry integrity for ${packageVersion} differs from the checked-main artifact.`);
    }
  }
  return { publishProjects, tagsToPush };
}

function validateOrdinaryRevision(base, checked) {
  const removedPlans = git([
    'diff',
    '--diff-filter=DR',
    '--name-only',
    base,
    checked,
    '--',
    '.nx/version-plans',
  ]).split('\n').filter(path => path.endsWith('.md'));
  if (removedPlans.length) {
    fail(`Version plans may only be consumed by a bound release candidate: ${removedPlans.join(', ')}`);
  }

  const manifests = git(['diff', '--name-only', base, checked, '--', 'libs/*/package.json'])
    .split('\n')
    .filter(Boolean);
  for (const path of manifests) {
    if (!gitMayFail(['cat-file', '-e', `${base}:${path}`]) || !gitMayFail(['cat-file', '-e', `${checked}:${path}`])) continue;
    const before = parseJson(git(['show', `${base}:${path}`], false), `${path} at ${base}`);
    const after = parseJson(git(['show', `${checked}:${path}`], false), `${path} at ${checked}`);
    if (before.version !== after.version) {
      fail(`${path} version changed outside a bound release candidate.`);
    }
  }
}

export function validateGeneratedTree({ source, candidate, metadata }) {
  validateMetadataShape(metadata);
  if (metadata.reviewedSource !== resolveCommit(source)) {
    fail('Release metadata reviewed source does not match the expected source.');
  }
  if (resolveCommit(`${candidate}^`) !== metadata.reviewedSource) {
    fail('Generated candidate must be a direct child of the reviewed source.');
  }

  const sourcePlanPaths = git(['ls-tree', '-r', '--name-only', source, '--', '.nx/version-plans'])
    .split('\n')
    .filter(path => path.endsWith('.md'));
  const recordedPlanPaths = metadata.plans.map(plan => plan.path);
  if (stableJson(sourcePlanPaths.sort()) !== stableJson([...recordedPlanPaths].sort())) {
    fail('Release metadata does not cover the exact reviewed version-plan set.');
  }

  const requested = new Map();
  for (const plan of metadata.plans) {
    const content = git(['show', `${source}:${plan.path}`], false);
    if (sha256(content) !== plan.sha256) fail(`Version-plan hash mismatch for ${plan.path}.`);
    for (const [project, bump] of parseVersionPlan(content, plan.path)) {
      requested.set(project, strongestBump(requested.get(project), bump));
    }
    if (gitMayFail(['cat-file', '-e', `${candidate}:${plan.path}`])) {
      fail(`Generated candidate still contains ${plan.path}.`);
    }
  }

  const releaseProjects = metadata.releases.map(release => release.project).sort();
  if (stableJson([...requested.keys()].sort()) !== stableJson(releaseProjects)) {
    fail('Release metadata package scope differs from the reviewed version plans.');
  }

  for (const release of metadata.releases) {
    const expected = releaseForProject({
      project: release.project,
      bump: requested.get(release.project),
      source,
      candidate,
      requireLocalTag: false,
    });
    if (stableJson(expected) !== stableJson(release)) {
      fail(`Forged or stale release metadata for ${release.project}.`);
    }
    const candidateChangelog = git(['show', `${candidate}:${release.changelogPath}`], false);
    const sourceHasChangelog = gitMayFail(['cat-file', '-e', `${source}:${release.changelogPath}`]);
    if (sourceHasChangelog && git(['show', `${source}:${release.changelogPath}`], false) === candidateChangelog) {
      fail(`${release.changelogPath} did not change in the generated candidate.`);
    }
  }

  const allowed = new Set([
    metadataPath,
    'pnpm-lock.yaml',
    ...metadata.plans.map(plan => plan.path),
    ...metadata.releases.flatMap(release => [release.manifestPath, release.changelogPath]),
  ]);
  const changed = git(['diff', '--name-only', source, candidate])
    .split('\n')
    .filter(Boolean);
  const unrelated = changed.filter(path => !allowed.has(path));
  if (unrelated.length) fail(`Generated candidate changes unrelated files: ${unrelated.join(', ')}`);
}

function releaseForProject({ project, bump, source, candidate, requireLocalTag = true }) {
  const projectPath = findProjectPath(project, candidate);
  const family = projectPath.split('/')[1];
  const manifestPath = `libs/${family}/package.json`;
  const changelogPath = `libs/${family}/CHANGELOG.md`;
  const before = parseJson(git(['show', `${source}:${manifestPath}`], false), `${manifestPath} at source`);
  const after = parseJson(git(['show', `${candidate}:${manifestPath}`], false), `${manifestPath} at candidate`);
  assertStableSemVer(before.version, `${manifestPath} source version`);
  assertStableSemVer(after.version, `${manifestPath} candidate version`);
  const expectedVersion = bumpVersion(before.version, bump);
  if (after.version !== expectedVersion) {
    fail(`${project} ${bump} plan resolves to ${expectedVersion}, not ${after.version}.`);
  }
  if (typeof after.name !== 'string' || !after.name.startsWith('@sneat/')) {
    fail(`${manifestPath} has no publishable @sneat package name.`);
  }
  const npmTag = `${project}-v${after.version}`;
  if (requireLocalTag && !git(['tag', '--points-at', candidate]).split('\n').includes(npmTag)) {
    fail(`Generated candidate is missing ${npmTag}.`);
  }
  return {
    family,
    project,
    package: after.name,
    manifestPath,
    changelogPath,
    version: after.version,
    npmTag,
    siblingGoTag: gitMayFail(['cat-file', '-e', `${candidate}:${family}/go.mod`])
      ? `${family}/v${after.version}`
      : null,
  };
}

function findProjectPath(project, revision) {
  const paths = git(['ls-tree', '-r', '--name-only', revision, '--', 'libs'])
    .split('\n')
    .filter(path => path.endsWith('/project.json'));
  const matches = paths.filter(path => {
    const value = parseJson(git(['show', `${revision}:${path}`], false), `${path} at ${revision}`);
    return value.name === project;
  });
  if (matches.length !== 1) fail(`Expected one project.json for ${project}, found ${matches.length}.`);
  return matches[0];
}

function parseVersionPlan(content, path) {
  const match = content.match(/^---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/);
  if (!match) fail(`${path} has no version-plan frontmatter.`);
  const entries = [];
  for (const line of match[1].split(/\r?\n/)) {
    if (!line.trim() || line.trimStart().startsWith('#')) continue;
    const field = line.match(/^([A-Za-z0-9_-]+):\s*(patch|minor|major)\s*$/);
    if (!field) fail(`${path} has unsupported version-plan field: ${line}`);
    entries.push([field[1], field[2]]);
  }
  if (!entries.length) fail(`${path} requests no package release.`);
  return entries;
}

function validateMetadataShape(value) {
  if (!value || value.schemaVersion !== 1 || !/^[0-9a-f]{40}$/.test(value.reviewedSource ?? '')) {
    fail('Release candidate metadata has an invalid schema or reviewed source.');
  }
  if (!Array.isArray(value.plans) || !value.plans.length || !Array.isArray(value.releases) || !value.releases.length) {
    fail('Release candidate metadata must include plans and releases.');
  }
  const planPaths = new Set();
  for (const plan of value.plans) {
    if (!/^\.nx\/version-plans\/[A-Za-z0-9._-]+\.md$/.test(plan.path ?? '') || !/^[0-9a-f]{64}$/.test(plan.sha256 ?? '')) {
      fail('Release candidate metadata contains an invalid plan receipt.');
    }
    if (planPaths.has(plan.path)) fail(`Duplicate version plan ${plan.path}.`);
    planPaths.add(plan.path);
  }
  const projects = new Set();
  const tags = new Set();
  for (const release of value.releases) {
    if (projects.has(release.project)) fail(`Duplicate release project ${release.project}.`);
    if (tags.has(release.npmTag)) fail(`Duplicate release tag ${release.npmTag}.`);
    projects.add(release.project);
    tags.add(release.npmTag);
  }
}

function readMetadata(revision) {
  return parseJson(git(['show', `${revision}:${metadataPath}`], false), `${metadataPath} at ${revision}`);
}

function sameFile(left, right, path) {
  const leftExists = gitMayFail(['cat-file', '-e', `${left}:${path}`]);
  const rightExists = gitMayFail(['cat-file', '-e', `${right}:${path}`]);
  if (leftExists !== rightExists) return false;
  if (!leftExists) return true;
  return git(['show', `${left}:${path}`], false) === git(['show', `${right}:${path}`], false);
}

export function bumpVersion(version, bump) {
  const [major, minor, patch] = version.split('.').map(Number);
  // Nx release keeps 0.x packages below 1.0: a "major" plan on 0.x bumps the
  // minor, and a "minor" plan bumps the patch (observed on main 2026-09-07:
  // budgetus-contract 0.2.6 + minor plan -> 0.2.7). Mirror that here so the
  // expected version matches what nx actually writes; a 1.0.0 of any contract
  // is an explicit decision, never a side effect of a plan.
  const effective = major === 0 ? (bump === 'major' ? 'minor' : bump === 'minor' ? 'patch' : bump) : bump;
  if (effective === 'major') return `${major + 1}.0.0`;
  if (effective === 'minor') return `${major}.${minor + 1}.0`;
  return `${major}.${minor}.${patch + 1}`;
}

function strongestBump(left, right) {
  const order = { patch: 1, minor: 2, major: 3 };
  return !left || order[right] > order[left] ? right : left;
}

function sha256(value) {
  return createHash('sha256').update(value).digest('hex');
}

function stableJson(value) {
  return JSON.stringify(value);
}

function parseJson(content, label) {
  try {
    return JSON.parse(content);
  } catch (error) {
    fail(`${label} is invalid JSON: ${error.message}`);
  }
}

function resolveCommit(revision) {
  return git(['rev-parse', `${revision}^{commit}`]);
}

function ensureAncestor(ancestor, descendant, message) {
  if (!gitMayFail(['merge-base', '--is-ancestor', ancestor, descendant])) fail(message);
}

function gitMayFail(args) {
  try {
    execFileSync('git', args, { stdio: 'ignore' });
    return true;
  } catch {
    return false;
  }
}

function git(args, trim = true) {
  try {
    const value = execFileSync('git', args, { encoding: 'utf8' });
    return trim ? value.trim() : value;
  } catch (error) {
    fail(`git ${args.join(' ')} failed: ${error.stderr?.toString().trim() || error.message}`);
  }
}

function assertStableSemVer(version, label) {
  if (!/^(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)$/.test(version ?? '')) {
    fail(`${label} is not stable SemVer: ${version}`);
  }
}

function fail(message) {
  throw new Error(message);
}

function writeOutput(value) {
  process.stdout.write(`${JSON.stringify(value)}\n`);
}

if (process.argv[1] && new URL(import.meta.url).pathname === process.argv[1]) {
  try {
    const command = process.argv[2];
    const source = process.env.SOURCE || 'HEAD^';
    const candidate = process.env.CANDIDATE || 'HEAD';
    if (command === 'prepare') {
      const metadata = createReleaseMetadata({ source, candidate });
      writeFileSync(metadataPath, `${JSON.stringify(metadata, null, 2)}\n`);
      writeOutput(metadata);
    } else if (command === 'verify-generated') {
      const metadata = JSON.parse(readFileSync(metadataPath, 'utf8'));
      validateGeneratedTree({ source, candidate, metadata });
      writeOutput(metadata);
    } else if (command === 'verify-final') {
      writeOutput(validateFinalCandidate({ checked: process.env.CHECKED || 'HEAD' }));
    } else if (command === 'verify-ci') {
      writeOutput(validateCiRevision({
        event: process.env.EVENT,
        revision: process.env.CHECKED || 'HEAD',
        base: process.env.BASE,
      }));
    } else if (command === 'classify') {
      writeOutput(classifyPublication({
        releases: JSON.parse(process.env.RELEASES || '[]'),
        checkedMain: process.env.CHECKED_MAIN,
        remoteTags: JSON.parse(process.env.REMOTE_TAGS || '{}'),
        npmIntegrities: JSON.parse(process.env.NPM_INTEGRITIES || '{}'),
        localIntegrities: JSON.parse(process.env.LOCAL_INTEGRITIES || '{}'),
      }));
    } else {
      fail('Usage: release-candidate.mjs prepare|verify-generated|verify-final|verify-ci|classify');
    }
  } catch (error) {
    console.error(`release candidate error: ${error.message}`);
    process.exit(1);
  }
}
