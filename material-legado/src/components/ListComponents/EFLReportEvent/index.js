import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_aux3_released.png';
import labelTabelaTituloImage from './assets/images/label_tabela_titulo.png';
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  container : {
    color: '#FFFFFF',
    overflow: 'auto',
    display: 'flex',
    flexDirection: 'column',
    justifyContent:"space-between"
  },
  top: {
    height: '26px',
    borderRadius: '20px 20px 0px 0px',
    backgroundColor: '#4864A0'
  },
  bottom:{
    height: '32px',
    borderRadius: '0px 0px 20px 20px',
    backgroundColor: '#4864A0'
  },
  title: {
    background: `url(${labelTabelaTituloImage}) transparent`,
    color: `#3C3C3C`,
    margin: 0,
    fontSize: '12px',
    height: '18px',
    fontWeight: 900,
    paddingLeft: '10px',
    lineHeight: '20px',
    borderLeft: '1px solid #7D8999',
    borderRight: '1px solid #7D8999',
  },
  ul:{
    margin: '0px',
    padding: '0px',
    flex: 1,
    borderLeft: '1px solid #7D8999',
    borderRight: '1px solid #7D8999',
    backgroundColor: '#B5C7DF',
  },
  li:{
    backgroundColor: '#A6B9D7',
    margin: '2px 0px 2px',
    paddingLeft: '10px',
    color: '#5064AA',
    fontSize: '12px',
    fontWeight: 'bold',
    borderTop: '1px solid #A6B9D7',
    borderBottom: '1px solid #A6B9D7'
  }
}

class EFLReportEvent extends JComponent {

  constructor(x, y) {
    super("EFLReportEvent", x, y);
    this.updateProperties();
    this.updateSchema();
  }

  updateProperties = () => {
    this.properties.label = "EFLReportEvent";
    this.properties.title = "EFLReportEvent";
    this.properties.name = "EFLReportEvent" + (+new Date());
    this.properties.height = 345;
    this.properties.width = 300;
    this.properties.option = "Exemplo1, Exemplo2, Exemplo3, Exemplo4";
  }

  updateSchema = () => {
    this.schema.basic.data.option = {
      label: 'Opções',
      key: 'option',
      type: 'string'
    };
    this.schema.basic.data.title = {
        label: 'Título',
        key: 'title',
        type: 'string'
      };
  }

  render() {
    return <div style={{ ...styles.container, ...this.getElementStyle(), backgroundSize: `${this.properties.width}px ${this.properties.height}px` }}>
        <div style={styles.top}></div>
        <p style={{ ...styles.title}}>{this.properties.title}</p>
        <ul style={styles.ul}>
            {
                this.properties.option.split(',').map((item, i) => (
                    <li style={{ ...styles.li}} key={i}>{item}</li>
                ))
            }
        </ul>
        <div style={styles.bottom}></div>
    </div>
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      return ['import br.gov.caixa.sispl.infra.ui.EFLReportEvent;', 
      'import br.gov.caixa.sispl.infra.ui.EFLReportEventListener;',
      'import br.gov.caixa.sispl.infra.ui.EFLPanel;',
      'import br.gov.caixa.sispl.infra.ui.EFLLabel;',
      'import br.gov.caixa.sispl.infra.skin.component.label.LabelSkin;'];
    },
    declaration: () => {
      return `    private EFLReportEvent ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tadd(get${pascalCase(this.properties.name)}());\n\taddOpcoes${pascalCase(this.properties.name)}();`;
    },
    method: () => {
      return `
    private EFLReportEvent get${pascalCase(this.properties.name)}(){
      if (${camelCase(this.properties.name)} == null) {
        ${camelCase(this.properties.name)} = new EFLReportEvent(${this.properties.width});
        ${camelCase(this.properties.name)}.setAutoscrolls(true);
        ${camelCase(this.properties.name)}.setLocation(${this.properties.left}, ${this.properties.top});
        ${camelCase(this.properties.name)}.setSize(${this.properties.width}, ${this.properties.height});
        ${camelCase(this.properties.name)}.addSelectListener(new ${pascalCase(this.properties.name)}ActionListener());
        EFLPanel cabecalho = ${camelCase(this.properties.name)}.addPanelCabecalho();
        addCabecalho${pascalCase(this.properties.name)}(cabecalho, "${this.properties.title}");
      }
      return ${camelCase(this.properties.name)};
    }

    private void addCabecalho${pascalCase(this.properties.name)}(EFLPanel cabecalho, String label) {
      EFLLabel labelCabecalho = new EFLLabel(LabelSkin.LABEL_REPORT_TITULO, label);
      labelCabecalho.setSize(cabecalho.getWidth(), cabecalho.getHeight());
      labelCabecalho.setLocation(10, 0);
      cabecalho.add(labelCabecalho);
    }

    private void addOpcoes${pascalCase(this.properties.name)}() {
    	// Remove os valores antes de incluir
    	get${pascalCase(this.properties.name)}().removeAll();
    	String[] opcoes = {${this.properties.option.split(",").map(i=>`"${i.trim()}"`).join(",")}};
    	for (String opcao : opcoes) {
    		EFLPanel panel = get${pascalCase(this.properties.name)}().addPanel();
    		panel.setSize(get${pascalCase(this.properties.name)}().getWidth(), 20);
    		EFLLabel label = new EFLLabel(LabelSkin.LABEL_TEXTO_AZUL, opcao);
    		label.setSize(panel.getWidth(), panel.getHeight());
    		label.setLocation(10, 0);
    		panel.add(label);
      }
    }

    private class ${pascalCase(this.properties.name)}ActionListener implements EFLReportEventListener {

      @Override
      public void onSelected(int index) {
        // TODO Auto-generated method stub
      }

    }
      `
    }
  })
}

const createImagePreview = () => {
  const dragIcon = document.createElement('img');
  dragIcon.src = image;
  dragIcon.width = 200;
  return dragIcon;
}

export default {
  id: 'EFLReportEvent',
  name: 'EFLReportEvent',
  component: EFLReportEvent,
  preview: {
    image: createImagePreview()
  }
}