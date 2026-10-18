import React from 'react';
import JComponent from '../JComponent';
import image from './assets/images/botao_aux3_released.png'
import {camelCase, pascalCase} from '../../../util/CodeUtil';

const styles = {
  background: `url(${image}) transparent`,
  backgroundRepeat: 'no-repeat',
  border: 'none',
  color: '#FFFFFF'
}

class EFLButton extends JComponent {

  constructor(x, y) {
    super("EFLButton", x, y);
    this.updateProperties();
    this.updateSchema();
  }

  updateProperties = () => {
    this.properties.label = "EFLButton";
    this.properties.name = "EFLButton" + (+new Date());
    this.properties.height = 60;
    this.properties.width = 110;
    this.properties.type = "CONFIRMAR";
  }

  events = {
    click: () => {
      if(this.properties.eventType === 'NAVEGAR' && this.context){
        const page = this.context.pages.find(p => p.id === +this.properties.eventPage);
        if(page){
          this.context.setSelectedPage(page);
        }
      }
    }
  }

  updateSchema = () => {
    this.schema.events = {
        label: 'Eventos',
        data: {
            eventType: {
                label: 'Tipo',
                key: 'eventType',
                event: 'click',
                eventFn: (context) => {
                  alert('Click')
                },
                type: 'select',
                options: ["NENHUM","NAVEGAR", "EXIBIR"]
            },
            eventValue: {
                label: 'Página',
                key: 'eventPage',
                type: 'select',
                optionsValues: (data) => {
                  if(data.pages){
                    const p = data.pages.map(p => p.id);
                    p.unshift('');
                    return p;
                  }
                  return [];
                },
                optionLabels:  (data) => {
                  if(data.pages){
                    const p = data.pages.map(p => p.name);
                    p.unshift('');
                    return p;
                  }
                  return [];
                },
                check: (data) => {
                  return data.properties.eventType === 'NAVEGAR';
                }
            },
        }
    };

    this.schema.basic.data.label = {
      label: 'Texto',
      key: 'label',
      type: 'string'
    };

    this.schema.basic.data.type = {
      label: 'Tipo',
      key: 'type',
      type: 'select',
      options: ["CONFIRMAR", "CANCELAR", "SIM", "NAO"]
    };
  }

  render(properties, context) {
    this.context = context;
    const props = {};
    if(this.properties.eventType !== 'NENHUM'){
      props.onClick = this.events.click.bind(this);
    }
    return <button {...props} style={{ ...styles, ...this.getElementStyle(), backgroundSize: `${this.properties.width}px ${this.properties.height}px` }}>{this.properties.label}</button>
  }

  /**
   * Renderiza o código do componente
   */
  code = (page) => ({
    imports: () => {
      return ['import br.gov.caixa.sispl.infra.ui.EFLButton;',
              'import java.awt.event.ActionEvent;',
              'import br.gov.caixa.sispl.infra.skin.component.button.ButtonSkin;',
              'import br.gov.caixa.sispl.infra.ui.SwingWorkerActionListener;'];
    },
    declaration: () => {
      return `    private EFLButton ${camelCase(this.properties.name)};`;
    },
    initialize: () => {
      return `\tadd(get${pascalCase(this.properties.name)}());`;
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

    private class ${pascalCase(this.properties.name)}ActionListener extends SwingWorkerActionListener {

      @Override
      public void actionPerformedInBackground(ActionEvent e) {
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
  id: 'EFLButton',
  name: 'EFLButton',
  component: EFLButton,
  preview: {
    image: createImagePreview()
  }
}