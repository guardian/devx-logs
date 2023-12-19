import type { Args } from "https://deno.land/std@0.200.0/flags/mod.ts";
import { parse } from "https://deno.land/std@0.200.0/flags/mod.ts";
import { getLink } from "./elk.ts";
import { parseFilters, removeUndefined } from "./transform.ts";

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

function main(inputArgs: string[]) {
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
  const link = getLink(space, filters, column);

  console.log(link);

  if (follow) {
    new Deno.Command("open", { args: [link] }).spawn();
  }
}

// Learn more at https://deno.land/manual/examples/module_metadata#concepts
if (import.meta.main) {
  main(Deno.args);
}
