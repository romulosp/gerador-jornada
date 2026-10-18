import React, { Component } from 'react';
import { AppConsumer } from '../../service/AppContext';
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";
import { Container, Item } from './styles';

export default class StructureTab extends Component {
    render() {
        return <AppConsumer>
            {
                data => (
                    <ul style={Container}>
                        {
                            data.selectedPage && data.selectedPage.components.map(component => {
                                const style = {};
                                if(data.selected && data.selected.id===component.id){
                                    style.backgroundColor = '#f39a00';
                                    style.color = '#FFFFFF';
                                }
                                return (
                                    <div key={component.id+"_"}>
                                        <ContextMenuTrigger id={component.id.toString()+"_"}>
                                            <li style={{...Item,...style}} key={component.id} onClick={() => {data.updateValue('selected', component)}}>
                                                {component.properties.name}
                                            </li>
                                        </ContextMenuTrigger>
                                        <ContextMenu id={component.id.toString()+"_"}>
                                            <MenuItem data={{ id: component.id }} onClick={() => data.removeComponent(component.id)}>
                                                Remover
                                            </MenuItem>
                                        </ContextMenu>
                                    </div>
                                )
                            })
                        }
                    </ul>
                )
            }
        </AppConsumer>;
    }
}
