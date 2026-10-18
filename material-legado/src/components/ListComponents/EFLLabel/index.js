import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_aux3_released.png';
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  border: 'none',
  fontWeight: 'bold',
  display: 'block'
}

class EFLLabel extends JComponent {

  constructor(x, y) {
    super("EFLLabel", x, y);
    this.updateProperties();
    this.updateSchema();
  }

  updateProperties = () => {
    this.properties.label = "EFLLabel";
    this.properties.name = "EFLLabel" + (+new Date());
    this.properties.height = 12;
    this.properties.width = 150;
    this.properties.fontSize = 12;
    this.properties.color = "#4864A0";
  }

  updateSchema = () => {
    this.schema.basic.data.label = {
      label: 'Texto',
      key: 'label',
      type: 'string'
    };

    this.schema.basic.data.fontSize = {
      label: 'Tamanho do Texto',
      key: 'fontSize',
      type: 'number',
      min: 6,
      max: 60
    };

    this.schema.basic.data.color = {
      label: 'Color',
      key: 'color',
      isStyle: true,
      renderStyle: (val) => val,
      type: 'color'
    };
  }

  render() {
    return <span style={{ ...styles, ...this.getElementStyle(), fontSize: `${this.properties.fontSize}px` }}>{this.properties.label}</span>
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      return ['import br.gov.caixa.sispl.infra.ui.EFLLabel;', 'import br.gov.caixa.sispl.infra.skin.component.label.LabelSkin;'];
    },
    declaration: () => {
      return `    private EFLLabel ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tadd(get${pascalCase(this.properties.name)}());`;
    },
    method: () => {
      return `
    public EFLLabel get${pascalCase(this.properties.name)}(){
      if (${camelCase(this.properties.name)} == null) {
        ${camelCase(this.properties.name)} = new EFLLabel(LabelSkin.LABEL_TEXTO_AZUL, "${this.properties.label}");
        ${camelCase(this.properties.name)}.setLocation(${this.properties.left}, ${this.properties.top});
        ${camelCase(this.properties.name)}.setSize(${this.properties.width}, ${this.properties.height});
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
  id: 'EFLLabel',
  name: 'EFLLabel',
  component: EFLLabel,
  preview: {
    image: createImagePreview()
  }
}