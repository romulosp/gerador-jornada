import React, { Component } from 'react';
import { AppProvider } from './service/AppContext'
import HomeScreen from './screens/HomeScreen'
import './styles.css';
import PageService from './service/PageService';
import ListComponents from './components/ListComponents';
import randomWords from 'random-words';
import uuidv4 from'uuid/v4';
 
export default class App extends Component {

  constructor(props) {
    super(props);
    this.state = {
      viewMode: "design",
      pages: [],
      components: ListComponents,
      selectedPage: null,
      selected: null,
    };
  }

  componentDidMount = () => {
    this.setState({
      pages: this._importPage(PageService.getPages()),
      selectedPage: PageService.getSelected()
    }, this.notifyUpdate);
  }
  

  updateValue = (key, val) => {
    this.setState({ [key]: val }, this.notifyUpdate);
  }

  updateSelectedPage = (val) => {
    this.setState(prev => {
        const selectedPage = val;
        const pages = prev.pages.filter(item => item.id !== val.id);
        pages.push(selectedPage);
        return {selectedPage, pages};
    }, this.notifyUpdate);
  }

  _importPage = (pages) => {
    return [...pages].map(p => {
      p.components = p.components.map(c => {
        const _c = this.state.components.find(i => i.name === c.name);
        const _instancia = new _c.component(0,0);
        _instancia.properties = c.properties;
        _instancia.id = c.id;
        _instancia.name = c.name;
        return _instancia;
      });
      return p;
    });
  }
  importPage = (pages) => {
    const newPages = this._importPage(pages);
    this.setState({
      selectedPage: newPages[0],
      pages: newPages
    }, this.notifyUpdate);
  }

  /**
   * Adiciona um componente a página corrente
   */
  addComponent = (component) => {
    this.setState(prevState => {
      let selectedPage = prevState.selectedPage;
      selectedPage.components = [...selectedPage.components, component];
      return { selectedPage, selected: component };
    }, this.notifyUpdate);
  }

  /**
   * Remove um elemento da página corrente
   */
  removeComponent = (id) => {
    let selectedPage = this.state.selectedPage;
    selectedPage.components = selectedPage.components.filter(element => element.id !== id);
    this.setState({
      selectedPage,
      selected: null
    }, this.notifyUpdate);
  }

  addPage = (newPage) => {
    const name = randomWords({ min: 3, max: 5, join: ' '});
    let page = {
      selected: true,
      showTitle: true,
      order: 0,
      title: name,
      name: name,
      components: []
    };
    if(newPage !== undefined){
      newPage = Object.assign({}, newPage);
      newPage.components = newPage.components.map(item => {
        const newItem = Object.assign({}, item);
        newItem.id = uuidv4();
        return newItem;
      });
      page = newPage;
    }
    page.id = +new Date();
    this.setState((oldProps) => {
      const pages = [...oldProps.pages].map(item => {
        item.selected = false;
        return item;
      });
      pages.push(page);

      console.log(pages);
      return {pages: this._importPage(pages), selectedPage: page, selected: null}
    }, this.notifyUpdate);
  }

  removePage = (pageId) => {
    this.setState((prev) => {
      const pages = prev.pages.filter(item => item.id !== pageId);

      return {pages, selectedPage: null, selected:null};
    }, this.notifyUpdate);
  }

  notifyUpdate = () => {
    PageService.setPages(this.state.pages);
  }

  setSelectedPage = (selectedPage) => {
    this.setState({selectedPage: selectedPage, selected: null})
  }

  render() {
    return (
      <AppProvider value={{
        removeComponent: this.removeComponent,
        components: this.state.components,
        addComponent: this.addComponent,
        viewMode: this.state.viewMode,
        updateValue: this.updateValue,
        pages: this.state.pages,
        selectedPage: this.state.selectedPage,
        selected: this.state.selected,
        importPage: this.importPage,
        addPage: this.addPage,
        notifyUpdate: this.notifyUpdate,
        setSelectedPage: this.setSelectedPage,
        removePage: this.removePage,
        updateSelectedPage: this.updateSelectedPage
      }}>
        <HomeScreen></HomeScreen>
      </AppProvider>
    );
  }
}
