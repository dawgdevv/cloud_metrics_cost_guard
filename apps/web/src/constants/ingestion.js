export const INITIAL_INGESTION_FORM = {
  account_id: "123456789012",
  days: 7,
  source: "synthetic",
  scenario: "gpu_spike",
};

export const INGESTION_SOURCE_STATUS = {
  synthetic: {
    label: "Demo",
    available: true,
    description: "Synthetic billing generator for local demos and testing.",
  },
  aws: {
    label: "AWS",
    available: true,
    description:
      "Fetches daily service-level spend from AWS Cost Explorer using your configured AWS credentials.",
  },
};

export const SCENARIO_LABELS = {
  gpu_spike: "GPU spike",
  storage_growth: "Storage growth",
  network_burst: "Network burst",
  normal: "Normal baseline",
};

export const SOURCE_OPTIONS = [
  { value: "synthetic", label: "Synthetic demo data" },
  { value: "aws", label: "AWS Cost Explorer" },
];