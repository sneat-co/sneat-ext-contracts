export function failedRequiredLegs(results) {
  const failed = ['NX_RESULT', 'TIER_COHERENCE_RESULT', 'DISCOVER_GO_RESULT'].filter(
    (name) => results[name] !== 'success',
  );
  if (!['success', 'skipped'].includes(results.GO_RESULT)) failed.push('GO_RESULT');
  return failed;
}

if (process.argv[1] === new URL(import.meta.url).pathname) {
  const failed = failedRequiredLegs(process.env);
  if (failed.length) {
    console.error(`Required CI legs did not pass: ${failed.join(', ')}`);
    process.exitCode = 1;
  }
}
