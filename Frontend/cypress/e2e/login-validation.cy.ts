describe("Login validation", () => {
  it("shows required field errors on empty submit", () => {
    cy.visit("/login");

    cy.contains("button", "Sign In").click();

    cy.contains("Email is required.").should("be.visible");
    cy.contains("Password is required.").should("be.visible");
  });

  it("shows email format error for invalid email", () => {
    cy.visit("/login");

    cy.get('input[placeholder="Email Address"]').type("invalid-email");
    cy.get('input[placeholder="Password"]').type("Abcd1234");
    cy.contains("button", "Sign In").click();

    cy.contains("Please enter a valid email address.").should("be.visible");
  });
});
