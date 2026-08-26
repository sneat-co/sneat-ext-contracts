/** A single contact reference tagged with the role it fills on a document. */
export interface IDocContactRef {
  readonly role: string;
  readonly contactID: string;
}

export interface IDocLinkageItemRef {
  readonly extID: string;
  readonly collection: string;
  readonly itemID: string;
}

/**
 * One role-tagged linkage edge a document implies once it can be persisted
 * through `facade4linkage` (`SetRelated`/`update_item_relationships`).
 *
 * Fable: persist as linkage edges (facade4linkage, role-tagged ItemRef) — do
 * not bake a bespoke member/role array. This type + `toLinkageEdgeDrafts`
 * exist so the mapping from "contact-role selections in the new-document
 * form" to "linkage edges" is defined and tested, ready for the day a
 * linkage client is wired into this app; nothing here is written to the
 * document DBO.
 */
export interface IDocLinkageEdgeDraft {
  readonly role: string;
  readonly from: IDocLinkageItemRef;
  readonly to: IDocLinkageItemRef;
}
