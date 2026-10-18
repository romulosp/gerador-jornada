import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_aux3_released.png'
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  ul:{
    listStyle : 'none',
    margin: '0px',
    padding: '0px',
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    justifyContent:"space-between"
  },
  li:{
    margin: '0px',
    paddingLeft: '0px',
    color: '#5064AA',
    fontSize: '12px',
    fontWeight: 'bold',
    display: 'flex'
  },
  texto: {
    border: 'none',
    fontWeight: 'bold',
    display: 'block',
    color: 'rgb(72, 100, 160)',
    fontSize: '16px'
  }
}

class EFLRadio extends JComponent {

  constructor(x, y) {
    super("EFLRadio", x, y);
    this.updateProperties();
    this.updateSchema();
  }

  updateProperties = () => {
    this.properties.label = "EFLRadio";
    this.properties.name = "EFLBtnGrp" + (+new Date());
    this.properties.height = 44;
    this.properties.width = 110;
    this.properties.option = "Exemplo1, Exemplo2";
  }

  updateSchema = () => {
    this.schema.basic.data.option = {
      label: 'Opções',
      key: 'option',
      type: 'string'
    };

    this.schema.basic.data.type = {
      label: 'Orientação',
      key: 'orientacao',
      type: 'select',
      options: ["VERTICAL", "HORIZONTAL"]
    };
  }

  render(properties, context) {
    this.context = context;
    return <ul style={{ ...styles.ul, flexDirection: properties.orientacao === 'HORIZONTAL' ? 'row' : 'column', width: `${this.properties.width}px`, height: `${this.properties.height}px`}}>
            {
                this.properties.option.split(',').map((item, i) => (
                    <li style={{ ...styles.li}} key={i}>
                          <label style={{...styles.texto}}>
                            <input type="radio" id={item} name={properties.name} value={item} />
                            {item}
                          </label>
                    </li>
                ))
            }
          </ul>
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      return ['import br.gov.caixa.sispl.infra.ui.EFLLabel;',
              'import br.gov.caixa.sispl.infra.ui.EFLButtonGroup;',
              'import br.gov.caixa.sispl.infra.ui.EFLPanel;'];
    },
    declaration: () => {
      return `    private EFLButtonGroup ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tget${pascalCase(this.properties.name)}();`;
    },
    method: () => {
      return `
    public EFLButton get${pascalCase(this.properties.name)}(){
      if (${camelCase(this.properties.name)} == null) {
        ${camelCase(this.properties.name)} = new EFLButton(uiManager, ButtonSkin.BOTAO_${this.properties.type}, "${this.properties.label}");
        ${camelCase(this.properties.name)}.setLocation(${this.properties.left}, ${this.properties.top});
        ${camelCase(this.properties.name)}.setSize(${this.properties.width}, ${this.properties.height});
        ${camelCase(this.properties.name)}.addActionListener(new ${pascalCase(this.properties.name)}ActionListener());
      }
      return ${camelCase(this.properties.name)};
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
  id: 'EFLRadio',
  name: 'EFLRadio',
  component: EFLRadio,
  preview: {
    image: createImagePreview()
  }
}