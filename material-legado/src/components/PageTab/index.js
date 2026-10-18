import React, { Component } from 'react';
import { AppConsumer } from '../../service/AppContext';
import { Container, Item, ButtonContainer, Button } from './styles';
import {MdAdd} from 'react-icons/md';
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";

export default class PageTab extends Component {
    render() {
        return <AppConsumer>
            {
                data => (
                    <>
                        <div style={ButtonContainer}>
                            <button style={Button} onClick={() => data.addPage()}>
                                <MdAdd />
                            </button>
                        </div>
                        <ul style={Container}>
                            {
                                data.pages.sort((a, b) => {
                                      const _a = a.order || a.id;
                                      const _b = b.order || b.id;
                                      if (_a < _b) {
                                        return -1;
                                      }
                                      if (_a > _b) {
                                        return 1;
                                      }
                                      return 0;
                                }).map(page => {
                                    const style = {};
                                    if(data.selectedPage && data.selectedPage.id===page.id){
                                        style.backgroundColor = '#f39a00';
                                        style.color = '#FFFFFF';
                                    }
                                    return <div key={page.id}>
                                        <ContextMenuTrigger id={page.id.toString()}>
                                            <li style={{...Item, ...style}} key={page.id} onClick={() => data.setSelectedPage(page)}>{page.name}</li>
                                        </ContextMenuTrigger>
                                        <ContextMenu id={page.id.toString()}>
                                            <MenuItem data={{ id: page.id }} onClick={() => data.removePage(page.id)}>
                                                Remover
                                            </MenuItem>
                                        </ContextMenu>
                                    </div>
                                })
                            }
                        </ul>
                    </>
                )
            }
        </AppConsumer>;
    }
}
