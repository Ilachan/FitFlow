describe("Register validation", () => {
  it("shows required field errors on empty submit", () => {
    cy.visit("/register");

    cy.contains("button", "Create Account").click();

    cy.contains("Full name is required.").should("be.visible");
    cy.contains("Email is required.").should("be.visible");
    cy.contains("Password is required.").should("be.visible");
  });

  it("validates password policy on weak password", () => {
    cy.visit("/register");

    cy.get('input[placeholder="Full Name"]').type("Test User");
    cy.get('input[placeholder="Email Address"]').type("test@example.com");
    cy.get('input[placeholder="Password"]').type("weak");
    cy.contains("button", "Create Account").click();

    cy.contains("Password must be 8+ chars").should("be.visible");
  });

  it("requires invitation code for manager registration", () => {
    cy.visit("/register");

    cy.contains("button", "Manager").click();
    cy.get('input[placeholder="Full Name"]').type("Manager Demo");
    cy.get('input[placeholder="Email Address"]').type("manager@example.com");
    cy.get('input[placeholder="Password"]').type("Abcd1234");
    cy.contains("button", "Create Account").click();

    cy.contains("Invitation code is required for managers.").should("be.visible");
  });
});
