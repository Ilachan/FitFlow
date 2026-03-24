import { describe, expect, it } from "vitest";
import { extractUserIdFromToken, roleIdToFrontendRole } from "./api";

const toBase64Url = (value: string) => {
  return btoa(value).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
};

describe("api helpers", () => {
  it("maps backend role id to frontend role", () => {
    expect(roleIdToFrontendRole(1)).toBe("student");
    expect(roleIdToFrontendRole(2)).toBe("manager");
    expect(roleIdToFrontendRole(3)).toBe("manager");
    expect(roleIdToFrontendRole(999)).toBe("student");
  });

  it("extracts user id from a valid JWT-like token payload", () => {
    const payload = { id: 42, role_id: 1 };
    const token = `header.${toBase64Url(JSON.stringify(payload))}.signature`;

    expect(extractUserIdFromToken(token)).toBe(42);
  });

  it("returns null for invalid token formats", () => {
    expect(extractUserIdFromToken("not-a-jwt")).toBeNull();
    expect(extractUserIdFromToken("header.invalid-base64.signature")).toBeNull();
  });

  it("returns null when payload has no numeric id", () => {
    const payload = { id: "42" };
    const token = `header.${toBase64Url(JSON.stringify(payload))}.signature`;

    expect(extractUserIdFromToken(token)).toBeNull();
  });
});
