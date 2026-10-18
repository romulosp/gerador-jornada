import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_released.png';
import image_aux1 from './assets/images/botao_aux1_released.png';
import imagePreviewTeclado from './assets/images/PreviewTeclado.png';
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-evenly'
  },
  button: {
    width: '52px',
    height: '52px',
    background: `url(${image}) transparent`,
    backgroundRepeat: 'no-repeat',
    border: 'none',
    color: '#FFFFFF',
    fontSize: '24px',
    fontWeight: 'bold'
  },
  largeButton: {
    width: '80px',
    height: '60px',
    background: `url(${image_aux1}) transparent`,
    backgroundRepeat: 'no-repeat',
    border: 'none',
    color: '#FFFFFF',
    fontSize: '12px',
    fontWeight: 'bold'
  }
}

class EFLVirtualKeyboardNumeric extends JComponent {

  constructor(x, y) {
    super("EFLVirtualKeyboardNumeric",x, y);
    this.updateProperties();
  }

  updateProperties = () => {
    this.properties.label = "EFLVirtualKeyboardNumeric";
    this.properties.name = "EFLVirtualKeyboardNumeric" + (+new Date());
    this.properties.height = 280;
    this.properties.width = 200;
  }

  render() {

    return (<div style={{ ...styles.container, ...this.getElementStyle()}}>
            <button style={styles.button}>7</button>
            <button style={styles.button}>8</button>
            <button style={styles.button}>9</button>
            <button style={styles.button}>4</button>
            <button style={styles.button}>5</button>
            <button style={styles.button}>6</button>
            <button style={styles.button}>1</button>
            <button style={styles.button}>2</button>
            <button style={styles.button}>3</button>
            <button style={styles.button}>&lt;</button>
            <button style={styles.button}>0</button>
            <button style={styles.button}>&gt;</button>
            <button style={styles.largeButton}>TAB</button>
            <button style={styles.largeButton}>LIMPAR</button>
          </div>)
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      return ['import br.gov.caixa.sispl.infra.ui.EFLVirtualKeyboardNumeric;'];
    },
    declaration: () => {
      return `    private EFLVirtualKeyboardNumeric ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tadd(get${pascalCase(this.properties.name)}());`;
    },
    method: () => {
      return `
    public EFLVirtualKeyboardNumeric get${pascalCase(this.properties.name)}(){
      if (${camelCase(this.properties.name)} == null) {
        ${camelCase(this.properties.name)} = new EFLVirtualKeyboardNumeric(uiManager);
        ${camelCase(this.properties.name)}.setLocation(${this.properties.left}, ${this.properties.top});
      }
      return ${camelCase(this.properties.name)};
    }
      `
    }
  })
}

const createImagePreview = () => {
  const dragIcon = document.createElement('img');
  dragIcon.src = imagePreviewTeclado;
  dragIcon.width = 200;
  return dragIcon;
}

export default {
  id: 'EFLVirtualKeyboardNumeric',
  name: 'EFLVirtualKeyboardNumeric',
  component: EFLVirtualKeyboardNumeric,
  preview: {
    image: createImagePreview()
  }
}