import { describe, expect, it } from "vitest";
import {
  evaluatePassword,
  formatPhoneNumberUS,
  getPasswordStrength,
  isAgeAtLeast,
  isPasswordValid,
  isValidAvatarUrl,
  isValidPhoneNumberUS,
  validateEmail,
} from "./validation";

describe("validation utilities", () => {
  it("validates email format", () => {
    expect(validateEmail("student@example.com")).toBe(true);
    expect(validateEmail("bad-email")).toBe(false);
  });

  it("evaluates password policy checks", () => {
    expect(evaluatePassword("Abcd1234")).toEqual({
      hasMinLength: true,
      hasUppercase: true,
      hasLowercase: true,
      hasNumber: true,
    });

    expect(isPasswordValid("Abcd1234")).toBe(true);
    expect(isPasswordValid("abcd1234")).toBe(false);
    expect(isPasswordValid("ABCD1234")).toBe(false);
  });

  it("returns password strength labels", () => {
    expect(getPasswordStrength("")).toBe("weak");
    expect(getPasswordStrength("abc")).toBe("weak");
    expect(getPasswordStrength("Abcdefgh")).toBe("medium");
    expect(getPasswordStrength("Abcdefg1")).toBe("strong");
  });

  it("formats and validates US phone numbers", () => {
    expect(formatPhoneNumberUS("1234567890")).toBe("(123) 456-7890");
    expect(isValidPhoneNumberUS("(123) 456-7890")).toBe(true);
    expect(isValidPhoneNumberUS("123-456-7890")).toBe(false);
  });

  it("validates avatar URLs", () => {
    expect(isValidAvatarUrl("https://example.com/avatar.png")).toBe(true);
    expect(isValidAvatarUrl("http://example.com/avatar.png")).toBe(true);
    expect(isValidAvatarUrl("example.com/avatar.png")).toBe(false);
  });

  it("checks minimum age", () => {
    expect(isAgeAtLeast("2000-01-01", 13)).toBe(true);
    expect(isAgeAtLeast("", 13)).toBe(true);

    const thisYear = new Date().getFullYear();
    expect(isAgeAtLeast(`${thisYear - 10}-01-01`, 13)).toBe(false);
  });
});
