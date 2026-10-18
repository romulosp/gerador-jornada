import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_aux3_released.png'
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  border: '1px solid #000000'
}

class EFLTextField extends JComponent {

  constructor(x, y) {
    super("EFLTextField", x, y);
    this.updateProperties();
    this.updateSchema();
  }

  updateProperties = () => {
    this.properties.label = "EFLTextField";
    this.properties.name = "EFLTextField" + (+new Date());
    this.properties.height = 25;
    this.properties.width = 150;
    this.properties.type = "Texto";
  }

  updateSchema = () => {
    this.schema.basic.data.type = {
      label: 'Tipo',
      key: 'type',
      type: 'select',
      options: ["Celular","Cep", "CPF", "Data","DDD", "Nis", "Número", "Texto"]
    };
  }

  render() {
    return <input type="text" style={{ ...styles, ...this.getElementStyle()}}></input>
  }

  getDocumentType = () =>{
    if(this.properties.type === 'Celular'){
      return `${camelCase(this.properties.name)}.setDocument(new EFLCelularDocument());`;
    }else if(this.properties.type === 'Cep'){
     return `${camelCase(this.properties.name)}.setDocument(new EFLCepDocument());`;
    }else if(this.properties.type === 'CPF'){
      return `${camelCase(this.properties.name)}.setDocument(new EFLCpfDocument());`;
    }else if(this.properties.type === 'Data'){
      return `${camelCase(this.properties.name)}.setDocument(new EFLDateDocument());`;
    }else if(this.properties.type === 'DDD'){
      return `${camelCase(this.properties.name)}.setDocument(new EFLDddDocument());`;
    }else if(this.properties.type === 'Nis'){
      return `${camelCase(this.properties.name)}.setDocument(new EFLNisDocument());`;
    }else if(this.properties.type === 'Número'){
     return `${camelCase(this.properties.name)}.setDocument(new EFLNumericDocument());`;
    }
    return '';
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      const imports = ['import br.gov.caixa.sispl.infra.ui.EFLTextField;', 'import br.gov.caixa.sispl.infra.skin.component.textfield.TextFieldSkin;'];
      if(this.properties.type === 'Celular'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLCelularDocument;');
      }else if(this.properties.type === 'Cep'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLCepDocument;');
      }else if(this.properties.type === 'CPF'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLCpfDocument;');
      }else if(this.properties.type === 'Data'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLDateDocument;');
      }else if(this.properties.type === 'DDD'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLDddDocument;');
      }else if(this.properties.type === 'Nis'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLNisDocument;');
      }else if(this.properties.type === 'Número'){
        imports.push('import br.gov.caixa.sispl.infra.ui.document.EFLNumericDocument;');
      }
      return imports;
    },
    declaration: () => {
      return `    private EFLTextField ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tadd(get${pascalCase(this.properties.name)}());`;
    },
    method: () => {
      return `
    public EFLTextField get${pascalCase(this.properties.name)}(){
      if (${camelCase(this.properties.name)} == null) {
        ${camelCase(this.properties.name)} = new EFLTextField(uiManager, TextFieldSkin.TEXT_LOGIN);
        ${camelCase(this.properties.name)}.setLocation(${this.properties.left}, ${this.properties.top});
        ${camelCase(this.properties.name)}.setSize(${this.properties.width}, ${this.properties.height});
        ${this.getDocumentType()}
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
  id: 'EFLTextField',
  name: 'EFLTextField',
  component: EFLTextField,
  preview: {
    image: createImagePreview()
  }
}