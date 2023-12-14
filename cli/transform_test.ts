import { assertEquals } from "https://deno.land/std@0.208.0/assert/assert_equals.ts";
import { parseFilters, removeUndefined } from "./transform.ts";

Deno.test("parseFilters", () => {
  const got = parseFilters(["stack=deploy", "stage=PROD", "app=riff-raff"]);
  const want = {
    stack: "deploy",
    stage: "PROD",
    app: "riff-raff",
  };
  assertEquals(got, want);
});

Deno.test("parseFilters without an = delimiter", () => {
  const got = parseFilters(["message"]);
  const want = {
    message: undefined,
  };
  assertEquals(got, want);
});

Deno.test("parseFilters without a value on the RHS of =", () => {
  const got = parseFilters(["name="]);
  const want = {
    name: "",
  };
  assertEquals(got, want);
});

Deno.test("removeUndefined", () => {
  const got = removeUndefined({
    stack: "deploy",
    stage: undefined,
    app: "riff-raff",
    team: "",
  });
  const want = {
    stack: "deploy",
    app: "riff-raff",
  };
  assertEquals(got, want);
});
