export interface AuthResponse {
  token: string;
  user: {
    id: number;
    username: string;
    name: string;
  };
}
