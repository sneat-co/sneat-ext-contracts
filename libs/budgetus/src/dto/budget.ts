export interface IBudgetusSpaceDbo {
  listGroups?: any[];
}

export interface IMoney {
  currency: string;
  value: number;
}

export type BudgetLineSource = 'asset-renewal' | 'happening' | 'gift';

export interface IBudgetLineItem {
  id: string;
  title: string;
  dateISO: string;
  amount: IMoney;
  source: BudgetLineSource;
  sourceRef?: string;
  targetAmount?: IMoney;
  isSurprise?: boolean;
}

export interface IBudgetMonthGroup {
  monthISO: string;
  total: IMoney;
  items: IBudgetLineItem[];
}

export interface IBudgetRollup {
  byMonth: IBudgetMonthGroup[];
  annualTotal: IMoney;
  mostExpensiveMonthISO: string;
}

export interface IBudgetOverridePatch {
  targetAmount?: IMoney;
  isSurprise?: boolean;
}
