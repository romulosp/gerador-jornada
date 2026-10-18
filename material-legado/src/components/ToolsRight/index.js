import React, { Component } from 'react';
import { Container } from './styles';
import DatGui, { DatBoolean, DatColor, DatNumber, DatString, DatFolder, DatSelect, DatPresets } from 'react-dat-gui';
import { AppConsumer } from '../../service/AppContext';
import 'react-dat-gui/dist/index.css';
import './dat.gui.css';

export default class ToolsRight extends Component {

    handleUpdate = newData => {
        const selected = this.selected;
        selected.properties = { ...selected.properties, ...newData }
        this.updateValue('selected', selected);
    }

    renderFolder = (title) => {
        return <DatFolder title={title}></DatFolder>
    }

    convertToArray = (data) => {
        if (!data) {
            return []
        }
        return Object.keys(data).map(key => data[key]);
    };

    handleUpdatePage = newData => {
        this.updateSelectedPage(newData);
    }

    render() {
        return (
            <AppConsumer>
                {
                    (dataConsumer) => {
                        this.updateValue = dataConsumer.updateValue;
                        this.selected = dataConsumer.selected;
                        this.updateSelectedPage = dataConsumer.updateSelectedPage;
                        return <div style={Container} className={this.props.className}>
                            {
                                !dataConsumer.selected && dataConsumer.selectedPage &&
                                <DatGui data={dataConsumer.selectedPage} onUpdate={this.handleUpdatePage}>
                                    <DatFolder title="Configuração da Página" closed={false} key={dataConsumer.selectedPage.id}>
                                        <DatString key="name" path="name" label="Nome" />
                                        <DatBoolean key="showTitle" path="showTitle" label="Exibir Título" />
                                        <DatString key="title" path="title" label="Título" />
                                        <DatNumber key="order" path="order" label="Ordem" />
                                    </DatFolder>
                                </DatGui>
                            }
                            {
                                dataConsumer.selected &&
                                <DatGui data={dataConsumer.selected.properties} onUpdate={this.handleUpdate}>
                                    {
                                        this.convertToArray(dataConsumer.selected.schema)
                                            .filter(item => this.convertToArray(item.data).length > 0)
                                            .map((item, i) => {
                                                if(item.check && item.check(this.selected.properties)){
                                                    return <div key={i}></div>
                                                }
                                                const children = this.convertToArray(item.data);
                                                return <DatFolder title={item.label} closed={false} key={item.label}>
                                                    {
                                                        children
                                                            .filter(el => ((typeof el.check === 'function' && el.check(this.selected) === true) || !el.check))
                                                            .map((el, index) => {
                                                                if (el.type === 'string') {
                                                                    return <DatString key={el.key} path={el.key} label={el.label} {...el} />
                                                                } else if (el.type === 'number') {
                                                                    return <DatNumber key={el.key} path={el.key} label={el.label} {...el} />
                                                                } else if (el.type === 'color') {
                                                                    return <DatColor key={el.key} path={el.key} label={el.label} {...el} />
                                                                } else if (el.type === 'boolean') {
                                                                    return <DatBoolean key={el.key} path={el.key} label={el.label} {...el} />
                                                                } else if (el.type === 'select') {
                                                                    return <DatSelect key={el.key} 
                                                                    optionLabels={el.optionLabels ? el.optionLabels(dataConsumer) : el.optionLabels} 
                                                                    options={typeof el.optionsValues === 'function' ? el.optionsValues(dataConsumer) : el.options} 
                                                                    path={el.key} label={el.label}/>
                                                                } else if (el.type === 'presets'){
                                                                    return <DatPresets label={el.label} {...el}/>
                                                                } else {
                                                                    return <span key={index}>Não Definido</span>
                                                                }
                                                            })
                                                    }
                                                </DatFolder>
                                            })
                                    }
                                </DatGui>
                            }
                        </div>
                    }
                }
            </AppConsumer>
        );
    }
}