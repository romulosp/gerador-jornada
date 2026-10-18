import randomWords from 'random-words';

class PageService {

    constructor() {
        const pages = localStorage.getItem('pagesSave');
        if(pages){
            try{
                this.pages = JSON.parse(pages).map(p => {
                    p.selected = false;
                    return p;
                });
            }catch(e){
                this.pages = [];
                this.pages.push(this._randomPage());
            }   
        }else{
            this.pages = [];
            this.pages.push(this._randomPage());
        }
    }

    _randomPage(){
        const name = randomWords({ min: 3, max: 5, join: ' '});
        return {
            id: +new Date(),
            selected: true,
            title: name,
            order: 0,
            name: name,
            showTitle: true,
            components: []
        }
    }

    getPages() {
        return this.pages;
    }

    setPages(pages) {
        this.pages = pages;
        localStorage.setItem('pagesSave', JSON.stringify(this.pages));
    }

    getSelected() {
        return this.getPages().find(page => page.selected);
    }
}

export default new PageService();