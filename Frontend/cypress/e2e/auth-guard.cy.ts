describe("Route protection", () => {
  it("redirects unauthenticated users from protected profile route to login", () => {
    cy.visit("/profile");

    cy.url().should("include", "/login");
    cy.contains("Sign In").should("be.visible");
  });
});
