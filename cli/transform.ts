/**
 * Turn an array of strings of form `key=value` into an object of form `{ key: value }`
 */
export function parseFilters(filter: string[]): Record<string, unknown> {
  return filter.reduce((acc, curr) => {
    const [key, value] = curr.split("=");
    return { ...acc, [key]: value };
  }, {});
}

/**
 * Remove keys from a `Record` whose value is falsy
 */
export function removeUndefined(
  obj: Record<string, string | undefined>,
): Record<string, string> {
  return Object.entries(obj).filter(([, value]) => !!value).reduce(
    (acc, [key, value]) => ({
      ...acc,
      [key]: value,
    }),
    {},
  );
}
