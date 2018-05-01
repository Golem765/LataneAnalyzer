export class User {
  id?: any;
  name?: string;
  link?: string;
  friends?: any[];
  type?: string;
  mentions?: { [key: string]: number };
  pressure?: number;
  influence?:number;
}
