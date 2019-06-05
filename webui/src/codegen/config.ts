export type GeneratorConfig = {
    packageName: string;
    funcPrefix: string;
    funcSuffix: string;
    sliceType: string;
    acceptLessThan: boolean;
    useTimSort: boolean;
};


export function symbol(name: string, notPrivateKeyWord: boolean, config: GeneratorConfig): string {
    const capitalName = name[0].toUpperCase() + name.slice(1);
    const combinedName = `${config.funcPrefix}${capitalName}${config.funcSuffix}`;
    if (!notPrivateKeyWord) {
        return combinedName[0].toLowerCase() + combinedName.slice(1);
    }
    return combinedName;
}
