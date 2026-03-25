export type PasswordCheckResult = {
  hasMinLength: boolean;
  hasUppercase: boolean;
  hasLowercase: boolean;
  hasNumber: boolean;
};

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const AVATAR_URL_REGEX = /^https?:\/\/.+/i;
const PHONE_US_REGEX = /^\(\d{3}\) \d{3}-\d{4}$/;

export const validateEmail = (email: string): boolean => {
  return EMAIL_REGEX.test(email.trim());
};

export const evaluatePassword = (password: string): PasswordCheckResult => {
  return {
    hasMinLength: password.length >= 8,
    hasUppercase: /[A-Z]/.test(password),
    hasLowercase: /[a-z]/.test(password),
    hasNumber: /\d/.test(password),
  };
};

export const isPasswordValid = (password: string): boolean => {
  const checks = evaluatePassword(password);
  return (
    checks.hasMinLength &&
    checks.hasUppercase &&
    checks.hasLowercase &&
    checks.hasNumber
  );
};

export const formatPhoneNumberUS = (input: string): string => {
  const normalized = input.replace(/\D/g, "").slice(0, 10);

  if (!normalized) {
    return "";
  }

  const area = normalized.slice(0, 3);
  const exchange = normalized.slice(3, 6);
  const line = normalized.slice(6, 10);

  if (normalized.length <= 3) {
    return `(${area}`;
  }

  if (normalized.length <= 6) {
    return `(${area}) ${exchange}`;
  }

  return `(${area}) ${line ? `${exchange}-${line}` : exchange}`;
};

export const isValidPhoneNumberUS = (phone: string): boolean => {
  if (!phone.trim()) {
    return true;
  }
  return PHONE_US_REGEX.test(phone.trim());
};

export const isValidAvatarUrl = (url: string): boolean => {
  if (!url.trim()) {
    return true;
  }
  return AVATAR_URL_REGEX.test(url.trim());
};

export const isAgeAtLeast = (dateOfBirth: string, minAge: number): boolean => {
  if (!dateOfBirth) {
    return true;
  }

  const birthDate = new Date(dateOfBirth);
  if (Number.isNaN(birthDate.getTime())) {
    return false;
  }

  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthDate.getDate())
  ) {
    age -= 1;
  }

  return age >= minAge;
};
