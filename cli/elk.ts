function escapeColon(str: string): string {
  return str.includes(":") ? `'${str}'` : str;
}
export function getLink(
  space: string,
  filters: Record<string, string>,
  columns: string[] = [],
): string {
  const kibanaFilters = Object.entries(filters).map(([key, value]) => {
    return `(query:(match_phrase:(${escapeColon(key)}:'${value}')))`;
  });

  // The `#/` at the end is important for Kibana to correctly parse the query string
  // The `URL` object moves this to the end of the string, which breaks the link.
  const base = `https://logs.gutools.co.uk/s/${space}/app/discover#/`;

  const query = {
    ...(kibanaFilters.length > 0 && {
      _g: `(filters:!(${kibanaFilters.join(",")}))`,
    }),
    ...(columns.length > 0 && {
      _a: `(columns:!(${columns.map(escapeColon).join(",")}))`,
    }),
  };

  const queryString = Object.entries(query)
    .map(([key, value]) => `${key}=${value}`)
    .join("&");

  return `${base}?${queryString}`;
}
