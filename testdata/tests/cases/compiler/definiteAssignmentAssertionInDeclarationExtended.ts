// @declaration: true
// @strict: true

export class DbObject {
    id!: string;
    name?: string;
    count: number = 0;
    private secret!: string;
    protected value!: number;
    static config!: boolean;
}

export interface IConfig {
    setting?: boolean; 
    optionalSetting?: string;
}