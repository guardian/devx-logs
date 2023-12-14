import { assertEquals } from "https://deno.land/std@0.208.0/assert/mod.ts";
import { getLink } from "./elk.ts";

// NOTE: Each of these URLs should be opened in a browser to verify that they work as expected.

Deno.test("getLink with simple input", () => {
  const got = getLink("devx", { app: "riff-raff", stage: "PROD" });
  const want =
    "https://logs.gutools.co.uk/s/devx/app/discover#/?_g=(filters:!((query:(match_phrase:('app':'riff-raff'))),(query:(match_phrase:('stage':'PROD')))))";
  assertEquals(got, want);
});

Deno.test("getLink with columns", () => {
  const got = getLink("devx", { app: "riff-raff", stage: "PROD" }, [
    "message",
    "level",
  ]);
  const want =
    "https://logs.gutools.co.uk/s/devx/app/discover#/?_g=(filters:!((query:(match_phrase:('app':'riff-raff'))),(query:(match_phrase:('stage':'PROD')))))&_a=(columns:!('message','level'))";
  assertEquals(got, want);
});

/*
Filters and columns with colon(:) input should get wrapped in single quotes(') so that Kibana can parse them correctly.
That is, gu:repo should become 'gu:repo'.
 */
Deno.test("getLink with colon(:) input", () => {
  const got = getLink("devx", {
    "gu:repo.keyword": "guardian/amigo",
    stage: "PROD",
  }, ["message", "gu:repo"]);
  const want =
    "https://logs.gutools.co.uk/s/devx/app/discover#/?_g=(filters:!((query:(match_phrase:('gu:repo.keyword':'guardian/amigo'))),(query:(match_phrase:('stage':'PROD')))))&_a=(columns:!('message','gu:repo'))";
  assertEquals(got, want);
});
