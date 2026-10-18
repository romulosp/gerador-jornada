import accents from 'remove-accents';
import _camelCase from 'camelcase';

export function pascalCase(text){
    if(!text){
        return '';
    }
    let data = accents.remove(text);
    return _camelCase(data, {pascalCase: true});
}

export function camelCase(text){
    if(!text){
        return '';
    }
    let data = accents.remove(text);
    return _camelCase(data);
}