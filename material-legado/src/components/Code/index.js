import React, { Component } from 'react';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { dark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Container } from './styles';
import { AppConsumer } from '../../service/AppContext';
import ejs from 'ejs';
import codeBaseString from './base';
import {camelCase, pascalCase} from '../../util/CodeUtil';

dark.hljs.margin = 0;
dark.hljs.flex = 1;
dark.hljs.backgroundColor = "rgb(22, 26, 32)";
dark.hljs.maxHeight = 'calc(100vh - 80px)';
dark.hljs.userSelect = 'text';

export default class Code extends Component {

  getCodeElements = (page) => {
    if(!page) return;
    const params = {};
    let elements = page.components || [];
    elements = elements.filter(el => typeof el.code === 'function');
    params.imports = new Set();
    params.imports.add("import br.gov.caixa.sispl.infra.ui.EFLPanel;");
    params.imports.add("import br.gov.caixa.sispl.infra.skin.component.panel.PanelSkin;");
    params.imports.add("import br.gov.caixa.sispl.infra.ui.EFLLabel;");
    if(elements){
      elements.map(el => el.code(page).imports()).forEach(item => {
        item.filter(item => item !=='').map(item => item.trim()).forEach(i => params.imports.add(i));
      });
    }
    params.imports= [...params.imports];
    params.initialize = elements.map(el => el.code(page).initialize()).filter(el => !!el && el!=='');
    params.declaration = elements.map(el => el.code(page).declaration()).filter(el => !!el && el!=='');
    params.method = elements.map(el => el.code(page).method()).filter(el => !!el && el!=='');
    return params;
  }

  render() {
    return (
      <AppConsumer>
        {
          props => {
            const page = props.selectedPage;
            if(!page){
              return <div style={Container} className={this.props.className}></div>
            }
            const codeElements = this.getCodeElements(page);
            const codeFormated = ejs.render(codeBaseString, { page, codeElements, camelCase, pascalCase});
            return (
              <div style={Container} className={this.props.className}>
                {page && 
                <SyntaxHighlighter language="java" style={dark}>
                  {codeFormated}
                </SyntaxHighlighter>
                }
              </div>
            )
          }
        }
      </AppConsumer>
    );
  }
}
