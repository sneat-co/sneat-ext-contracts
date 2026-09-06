export function failedRequiredLegs(results) {
  const failed = ['NX_RESULT', 'TIER_COHERENCE_RESULT', 'DISCOVER_GO_RESULT'].filter(
    (name) => results[name] !== 'success',
  );
  let goDirs;
  try {
    goDirs = JSON.parse(results.GO_DIRS);
  } catch {
    failed.push('GO_DIRS');
  }
  if (!Array.isArray(goDirs) || goDirs.some((dir) => typeof dir !== 'string' || !dir)) {
    if (!failed.includes('GO_DIRS')) failed.push('GO_DIRS');
  }
  const validGoResult =
    results.GO_RESULT === 'success' ||
    (results.GO_RESULT === 'skipped' && Array.isArray(goDirs) && goDirs.length === 0);
  if (!validGoResult) failed.push('GO_RESULT');
  return failed;
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  const failed = failedRequiredLegs(process.env);
  if (failed.length) {
    console.error(`Required CI legs did not pass: ${failed.join(', ')}`);
    process.exitCode = 1;
  }
}
import { fileURLToPath } from 'node:url';
