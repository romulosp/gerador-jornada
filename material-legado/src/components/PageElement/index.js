import React, { Component } from 'react';
import { Container } from './styles';
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";
export default class PageElement extends Component {

  constructor(props) {
    super(props);
    this.state = {
      hover: false
    }
  }

  onMouseOver = () => {
    if (!this.state.hover) {
      this.setState({ hover: true })
    }
  }

  onMouseLeave = () => {
    if (this.state.hover) {
      this.setState({ hover: false })
    }
  }

  render() {
    const { element, removeComponent, onMouseDown, onMouseUp, selected } = this.props;
    const properties = element.properties;
    const renderElement = element.render(properties);
    const selectedElementStyle = selected && selected.id === element.id ? { outline: '5px dotted #4864A0' } : {};
    return (
      <React.Fragment>
        <ContextMenuTrigger id={element.id.toString()} holdToDisplay={-1}>
          <div onClick={() => this.props.onClick(element)} onMouseLeave={this.onMouseLeave} onMouseOver={this.onMouseOver} style={{ ...Container, ...element.getContainerStyle(), ...selectedElementStyle }}>
            <div onMouseOver={this.onMouseOver} style={{ position: 'absolute', width: '100%', height: '100%' }} onMouseDown={() => onMouseDown(element)} onMouseUp={onMouseUp}>

            </div>
            {
              renderElement
            }
          </div>
        </ContextMenuTrigger>
        <ContextMenu id={element.id.toString()}>
          <MenuItem data={{ id: element.id }} onClick={(i, item) => removeComponent(item.id)}>
            Remover
            </MenuItem>
        </ContextMenu>
      </React.Fragment>
    );
  }
}