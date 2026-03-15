export function formatNumber(value) {
  if (value === null || value === undefined || value === "") {
    return "A definir";
  }

  return new Intl.NumberFormat("fr-FR").format(Number(value));
}

export function formatCurrency(value) {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
  }).format(Number(value ?? 0));
}

export function formatDate(value) {
  if (!value) {
    return "-";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric"
  }).format(date);
}
