export interface ContractGeneratorSchema {
  /**
   * Family name — lowercase, kebab-case (letters, digits, single hyphens,
   * starting with a letter — e.g. "taxus", "kids-club"). Becomes the Nx
   * project "<family>-contract" and npm package
   * "@sneat/extension-<family>-contract".
   */
  family: string;

  /**
   * Also scaffold <family>/go.mod (module
   * github.com/sneat-co/sneat-ext-contracts/<family>) with a doc.go stub,
   * and add "./<family>" to go.work's use block.
   */
  go?: boolean;
}
