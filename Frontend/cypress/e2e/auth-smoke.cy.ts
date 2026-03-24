describe("Auth smoke flow", () => {
  it("allows typing into login form and toggling password visibility", () => {
    cy.visit("/login");

    cy.get('input[placeholder="Email Address"]').type("student@example.com");
    cy.get('input[placeholder="Password"]').type("Abcd1234");

    cy.contains("button", "Show").click();
    cy.get('input[placeholder="Password"]').should("have.attr", "type", "text");

    cy.contains("button", "Hide").click();
    cy.get('input[placeholder="Password"]').should("have.attr", "type", "password");
  });
});
