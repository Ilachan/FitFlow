describe("Register validation", () => {
  it("shows required field errors on empty submit", () => {
    cy.visit("/register");

    cy.contains("button", "Create Account").click();

    cy.contains("Full name is required.").should("be.visible");
    cy.contains("Email is required.").should("be.visible");
    cy.contains("Password is required.").should("be.visible");
  });

  it("updates password strength indicator while typing", () => {
    cy.visit("/register");

    cy.contains("weak").should("be.visible");
    cy.get('input[placeholder="Password"]').type("Abcdefg1");
    cy.contains("strong").should("be.visible");
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
