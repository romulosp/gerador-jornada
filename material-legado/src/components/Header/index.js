import React, { Component } from 'react';
import { Container, Title, ButtonLabel, RightButtonContainer, Button, ButtonActive, LeftButtonContainer } from './styles';
import { AppConsumer } from '../../service/AppContext';
import html2canvas from 'html2canvas';
import {
    MdScreenShare,
    MdFormatAlignCenter,
    MdVerticalAlignCenter,
    MdFormatAlignLeft,
    MdFormatAlignRight,
    MdVerticalAlignBottom,
    MdVerticalAlignTop,
    MdSave
} from 'react-icons/md';

import { FaFileImport } from 'react-icons/fa';

export default class Header extends Component {


    onScreenshot = () => {
        html2canvas(document.querySelector("#page-preview")).then(canvas => {
            var link = document.createElement("a");
            link.setAttribute('download', `${+new Date()}.png`);
            link.setAttribute('href', canvas.toDataURL("image/png").replace("image/png", "image/octet-stream"));
            link.click();
        });
    }

    onCenter = (rootProps) => {
        const element = rootProps.selected;
        const widthElement = element.properties.width;
        const left = 400 - parseInt((widthElement / 2), 10);
        element.properties = { ...element.properties, left };
        rootProps.updateValue('selected', element);
    }

    onCenterVertical = (rootProps) => {
        const element = rootProps.selected;
        const heightElement = element.properties.height;
        const top = 300 - parseInt((heightElement / 2), 10);
        element.properties = { ...element.properties, top };
        rootProps.updateValue('selected', element);
    }

    onAlignLeft = (rootProps) => {
        const element = rootProps.selected;
        element.properties = { ...element.properties, left: 25 };
        rootProps.updateValue('selected', element);
    }

    onAlignRight = (rootProps) => {
        const element = rootProps.selected;
        const widthElement = element.properties.width;
        element.properties = { ...element.properties, left: (800 - widthElement - 25) };
        rootProps.updateValue('selected', element);
    }

    onAlignBottom = (rootProps) => {
        const element = rootProps.selected;
        const heightElement = element.properties.height;
        element.properties = { ...element.properties, top: (600 - heightElement - 25) };
        rootProps.updateValue('selected', element);
    }

    onAlignTop = (rootProps) => {
        const element = rootProps.selected;
        element.properties = { ...element.properties, top: (128 + 25) };
        rootProps.updateValue('selected', element);
    }

    onSave = (rootProps) => {
        const newPages = rootProps.pages.map(p => {
            const page = Object.assign({}, p);
            page.components = page.components.map(item => item.toJSON());
            return page;
        });
        const data = JSON.stringify(newPages);
        var file = new Blob([data], { type: "application/json" });
        var link = document.createElement("a");
        link.setAttribute('download', `${+new Date()}.page`);
        link.setAttribute('href', URL.createObjectURL(file));
        link.click();
    }

    onImport = (rootProps) => {
        const reader = new FileReader();
        reader.onload = () => {
            try {
                const text = reader.result;
                if (text) {
                    const data = JSON.parse(text);
                    rootProps.importPage(data);
                }
            } catch (e) {
                console.log("Ocorreu um erro")
            }
        }
        reader.readAsText(document.getElementById('selecao-arquivo').files[0]);
    }

    render() {
        return (
            <div style={Container} className={this.props.className}>
                <AppConsumer>
                    {props => (
                        <React.Fragment>
                            <div style={LeftButtonContainer}>
                                <input id="selecao-arquivo" style={{ display: 'none' }} type="file" onChange={() => this.onImport(props)} />
                                <div style={{ ...Button, display: props.pages && props.pages.length > 0 ? 'flex' : 'none' }} onClick={() => this.onSave(props)}>
                                    <MdSave />
                                    <span style={ButtonLabel}>Salvar</span>
                                </div>
                                <label style={{ ...Button }} htmlFor='selecao-arquivo'>
                                    <FaFileImport />
                                    <span style={ButtonLabel}>Importar</span>
                                </label>
                                <div style={{ ...Button, display: props.viewMode === 'design' && props.selectedPage ? 'flex' : 'none' }} onClick={this.onScreenshot}>
                                    <MdScreenShare />
                                    <span style={ButtonLabel}>Screenshot</span>
                                </div>
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onCenter(props)}>
                                    <MdFormatAlignCenter /></div>}
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onCenterVertical(props)}><MdVerticalAlignCenter /></div>}
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onAlignLeft(props)}><MdFormatAlignLeft /></div>}
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onAlignRight(props)}><MdFormatAlignRight /></div>}
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onAlignBottom(props)}><MdVerticalAlignBottom /></div>}
                                {props.viewMode === 'design' && <div style={{ ...Button, display: props.selected ? 'block' : 'none' }} onClick={() => this.onAlignTop(props)}><MdVerticalAlignTop /></div>}
                            </div>
                            <div style={Title}>

                            </div>
                            <div style={RightButtonContainer}>
                                <div style={props.viewMode === 'view' ? ButtonActive : Button} onClick={() => props.updateValue('viewMode', 'view')}>VIEW</div>
                                <div style={props.viewMode === 'design' ? ButtonActive : Button} onClick={() => props.updateValue('viewMode', 'design')}>DESIGN</div>
                                <div style={props.viewMode === 'code' ? ButtonActive : Button} onClick={() => props.updateValue('viewMode', 'code')}>CODE</div>
                            </div>
                        </React.Fragment>
                    )}
                </AppConsumer>
            </div>
        );
    }
}