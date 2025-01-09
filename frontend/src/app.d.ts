// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
  // interface Locals {}
  // interface PageData {}
  // interface Error {}
  // interface Platform {}
}

interface DataSet {
  label: string;
  data: number[];
  xAxis: string[];
  borderColor: string;
  backgroundColor: string;
}

interface LiquidityRecord {
  timestamp: string;
  bid_ask_spread: number;
  volume: number;
  bid_price: number;
}

interface LiquidityReport {
  asset_type: string;
  total_records: number;
  historical_records: number;
  prediction_records: number;
  current_high_risk_count: number;
  current_moderate_risk_count: number;
  predicted_high_risk_count: number;
  predicted_moderate_risk_count: number;
}
