import {
  IAssetDbo,
  IExpiringDocumentsResponse,
  IReQuoteBundle,
} from './index';

// Minimal valid asset dbo for building list items in the tests.
const asset = (id: string, category: IAssetDbo['category']): IAssetDbo =>
  ({
    id,
    name: id,
    category,
    status: 'active',
    visibility: 'private',
  }) as IAssetDbo;

describe('ReQuoter read shapes', () => {
  it('IReQuoteBundle round-trips vehicle + policies + licences', () => {
    const bundle: IReQuoteBundle = {
      vehicle: { id: 'veh1', asset: asset('veh1', 'vehicles') },
      policies: [{ id: 'pol1', asset: asset('pol1', 'document') }],
      licences: [{ id: 'dl1', asset: asset('dl1', 'document') }],
    };
    const roundTripped: IReQuoteBundle = JSON.parse(JSON.stringify(bundle));
    expect(roundTripped).toEqual(bundle);
    expect(roundTripped.vehicle?.id).toBe('veh1');
    expect(roundTripped.policies[0].id).toBe('pol1');
    expect(roundTripped.licences[0].id).toBe('dl1');
  });

  it('IReQuoteBundle vehicle is optional (unknown vehicle)', () => {
    const bundle: IReQuoteBundle = { policies: [], licences: [] };
    expect(bundle.vehicle).toBeUndefined();
  });

  it('IExpiringDocumentsResponse carries docType + expiresOn', () => {
    const resp: IExpiringDocumentsResponse = {
      documents: [
        {
          id: 'pol1',
          docType: 'insurance',
          expiresOn: '2027-01-01',
          asset: asset('pol1', 'document'),
        },
      ],
    };
    expect(resp.documents[0].docType).toBe('insurance');
    expect(resp.documents[0].expiresOn).toBe('2027-01-01');
  });
});
