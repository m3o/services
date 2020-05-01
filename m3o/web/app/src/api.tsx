import axios from 'axios';

const BaseURL = 'https://api.micro.mu/m3o/'

export async function Call(path: string, params?: any): Promise<any> {
  return axios.post(BaseURL + path, params, { withCredentials: true });
}

export interface Project {
  id?: string;
  name: string;
  description: string;
  repository?: string;
  api_domain?: string;
  web_domain?: string;
}

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  profile_picture_url?: string;
  email: string;
  roles: string[];
}

export interface EnvVar {
  id?: string;
  key: string;
  value: string;
  service: string;
  secret?: boolean;
}
