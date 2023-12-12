import type { Args } from "https://deno.land/std@0.200.0/flags/mod.ts";
import { parse } from "https://deno.land/std@0.200.0/flags/mod.ts";
import { open } from "https://deno.land/x/open@v0.0.6/index.ts";

export function getLink(
  space: string,
  filters: Record<string, string>,
  columns: string[],
): string {
  const kibanaFilters = Object.entries(filters).map(([key, value]) => {
    return `(query:(match_phrase:(${key}:'${value}')))`;
  });

  // The `#/` at the end is important for Kibana to correctly parse the query string
  // The `URL` object moves this to the end of the string, which breaks the link.
  const base = `https://logs.gutools.co.uk/s/${space}/app/discover#/`;

  const query = {
    ...(kibanaFilters.length > 0 && {
      _g: `(filters:!(${kibanaFilters.join(",")}))`,
    }),
    ...(columns.length > 0 && {
      _a: `(columns:!(${columns.join(",")}))`,
    }),
  };

  const queryString = Object.entries(query)
    .map(([key, value]) => `${key}=${value}`)
    .join("&");

  return `${base}?${queryString}`;
}

function parseArguments(args: string[]): Args {
  return parse(args, {
    boolean: ["follow"],
    negatable: ["follow"],
    string: ["space", "stack", "stage", "app"],
    collect: ["column", "filter"],
    stopEarly: false,
    "--": true,
    default: {
      follow: true,
      column: ["message"],
      space: "default",
      filter: [],
    },
  });
}

function escapeColon(str: string): string {
  return str.includes(":") ? `'${str}'` : str;
}

function parseFilters(filter: string[]): Record<string, string> {
  return filter.reduce((acc, curr) => {
    const [key, value] = curr.split("=");
    return { ...acc, [escapeColon(key)]: value };
  }, {});
}

function removeUndefined(
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

function printHelp(): void {
  console.log(`Usage: devx-logs [OPTIONS...]`);
  console.log("\nOptional flags:");
  console.log("  --help              Display this help and exit");
  console.log("  --space             The Kibana space to use");
  console.log("  --stack             The stack tag to filter by");
  console.log("  --stage             The stage tag to filter by");
  console.log("  --app               The app tag to filter by");
  console.log(
    "  --column            Which columns to display. Multiple: true. Default: 'message'",
  );
  console.log(
    "  --filter            Additional filters to apply. Multiple: true. Format: key=value",
  );
  console.log("  --no-follow         Don't open the link in the browser");
  console.log("\nExample:");
  console.log(
    "  devx-logs --space devx --stack deploy --stage PROD --app riff-raff",
  );
  console.log("\nAdvanced example:");
  console.log(
    "  devx-logs --space devx --stack deploy --stage PROD --app riff-raff --filter level=INFO --filter region=eu-west-1 --column message --column logger_name",
  );
}

async function main(inputArgs: string[]) {
  const args = parseArguments(inputArgs);

  if (args.help) {
    printHelp();
    Deno.exit(0);
  }

  const { space, stack, stage, app, column, filter, follow } = args;

  const mergedFilters: Record<string, string | undefined> = {
    ...parseFilters(filter),
    "stack.keyword": stack,
    "stage.keyword": stage,
    "app.keyword": app,
  };

  const filters = removeUndefined(mergedFilters);
  const link = getLink(space, filters, column.map(escapeColon));

  console.log(link);

  if (follow) {
    await open(link);
  }
}

await main(Deno.args);
