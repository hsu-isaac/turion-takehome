export const anomalyTypeDisplayNames: Record<string, string> = {
  high_temperature: "High Temperature",
  low_temperature: "Low Temperature",
  low_battery: "Low Battery",
  low_altitude: "Low Altitude",
  weak_signal: "Weak Signal",
};

export const getAnomalyDisplayName = (anomalyType: string): string => {
  return anomalyTypeDisplayNames[anomalyType] || anomalyType;
};
