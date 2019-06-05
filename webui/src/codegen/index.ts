import { GeneratorConfig } from './config';
import { generate } from './slices-template';
import { generateComparable } from './slices-comparable-template';
import { generateTimsort } from './slices-timsort-template';
import { generateComparableTimsort } from './slices-comparable-timsort-template';

export { GeneratorConfig } from './config';

export function generateSourceCode(config: GeneratorConfig): string {
    let source: string;
    if (config.acceptLessThan) {
        if (config.useTimSort) {
            source = generateTimsort(config);
        } else {
            source = generate(config);
        }
    } else {
        if (config.useTimSort) {
            source = generateComparableTimsort(config);
        } else {
            source = generateComparable(config);
        }
    }
    return simpleGoFormat(source);
}

function simpleGoFormat(src: string): string {
    const [_, ...lines] = src.split('\n');

    const result = lines.map((line: string) => {
        line = line.slice(4);
        const match = line.match(/^( +)/);
        if (match) {
            const indent = Math.floor(match[1].length / 4);
            const space = new Array(indent + 1).join('    ');
            const tab = new Array(indent + 1).join('\t');
            line = line.replace(space, tab);
        }
        return line;
    });

    return result.join('\n');
}
