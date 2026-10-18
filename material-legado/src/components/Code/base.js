const code = `
package br.gov.caixa.sispl.efl.uc.<%=pascalCase(page.name).toLowerCase()%>.ui.pnl;

<%- codeElements.imports.join("\\n"); %>

public class Pnl<%=pascalCase(page.name)%> extends EFLPanel {

    private static final long serialVersionUID = 1L;
    
<%- codeElements.declaration.join("\\n"); %>
<% if (page.showTitle) { %>
    private EFLPanel pnlTituloPrincipal = null;
    private EFLLabel labelTituloPrincipal = null;
<% } %>

    /** Controle de interface do caso de uso */
    protected UIControle ui; // TODO incluir controle

    public Pnl<%-pascalCase(page.name)%>(UIControle aUi) {
        super (PanelSkin.PANEL_EFL);
        this.ui = aUi;
        initialize();
    }
    
    private void initialize() {
    <%- codeElements.initialize.join("\\n"); %>
    <% if (page.showTitle) { %>
        add(getPnlTituloPrincipal());
    <% } %>
    }

    <% if (page.showTitle) { %>
    private EFLPanel getPnlTituloPrincipal() {
        if (pnlTituloPrincipal == null) {
            pnlTituloPrincipal = new EFLPanel(PanelSkin.PANEL_FAIXA_TITULO_PRINCIPAL);
            pnlTituloPrincipal.setLocation(0, 100);
            pnlTituloPrincipal.setSize(800, 28);
            pnlTituloPrincipal.setVisible(true);
            pnlTituloPrincipal.add(getLabelTituloPrincipal());
        }
        return pnlTituloPrincipal;
    }
    
    private EFLLabel getLabelTituloPrincipal() {
        if (labelTituloPrincipal == null) {
            labelTituloPrincipal = new EFLLabel(LabelSkin.LABEL_FAIXA_TITULO_PRINCIPAL, "<%-page.name%>", EFLLabel.RIGHT);
            labelTituloPrincipal.setLocation(0, 0);
            labelTituloPrincipal.setSize(750, 28);
        }
        return labelTituloPrincipal;
    }
    <% } %>

<%- codeElements.method.join("\\n"); %>
}
`;
export default code;