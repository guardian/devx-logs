/**
 * Wrap a string in single quotes so Kibana can parse it correctly
 */
function wrapString(str: string): string {
  return `'${str}'`;
}

function matchPhraseQuery(field: string, value: string) {
  return `(query:(match_phrase:(${wrapString(field)}:${wrapString(value)})))`;
}

function oneOfQuery(field: string, values: string[]) {
  const wrappedField = wrapString(field);

  const metaPart = `meta:(field:${wrappedField},key:${wrappedField},params:!(${
    values.map(wrapString).join(",")
  }),type:phrases,value:!(${values.map(wrapString).join(",")}))`;

  const queryPart = `query:(bool:(minimum_should_match:1,should:!(${
    values.map((value) =>
      `(match_phrase:(${wrappedField}:${wrapString(value)}))`
    )
  })))`;

  return `(${metaPart},${queryPart})`;
}

export function getLink(
  space: string,
  apps: string[] = [],
  phraseFilters: Record<string, string>,
  columns: string[] = [],
): string {
  const phraseQueries = Object.entries(phraseFilters).map(([key, value]) =>
    matchPhraseQuery(key, value)
  );

  const appQuery = apps.length > 0 ? [oneOfQuery("app.keyword", apps)] : [];

  const filters = [...appQuery, ...phraseQueries];

  // The `#/` at the end is important for Kibana to correctly parse the query string
  // The `URL` object moves this to the end of the string, which breaks the link.
  const base = `https://logs.gutools.co.uk/s/${space}/app/discover#/`;

  // TODO use https://github.com/Nanonid/rison to generate the query string
  const query = {
    ...(filters.length > 0 && {
      _g: `(filters:!(${filters.join(",")}))`,
    }),
    ...(columns.length > 0 && {
      _a: `(columns:!(${columns.map(wrapString).join(",")}))`,
    }),
  };

  const queryString = Object.entries(query)
    .map(([key, value]) => `${key}=${value}`)
    .join("&");

  return `${base}?${queryString}`;
}
