import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import axios from 'axios';
import {
  Container,
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  CircularProgress,
} from '@mui/material';
import type { AuthResponse } from '../types/user';
import { useNavigate } from 'react-router-dom';

interface FormData {
  name: string;
  username: string;
  password: string;
}

export function RegisterPage() {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    username: '',
    password: '',
  });

  const navigate = useNavigate();

  const mutation = useMutation<AuthResponse, Error, FormData>({
    mutationFn: async newUser => {
      const response = await axios.post(
        'http://localhost:3000/api/register',
        newUser
      );
      console.log(response);
      return response.data;
    },
    onSuccess: data => {
      localStorage.setItem('token', data.token);
      navigate('/dashboard', { replace: true });
    },
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (
      formData.name === '' ||
      formData.username === '' ||
      formData.password === ''
    ) {
      return;
    }
    mutation.mutate(formData);
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper
          elevation={3}
          sx={{ padding: 4, width: '100%', borderRadius: 2 }}
        >
          <Typography component="h1" variant="h5" align="center" gutterBottom>
            Create Account
          </Typography>

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <TextField
              margin="normal"
              required
              fullWidth
              id="name"
              label="Full Name"
              name="name"
              autoComplete="name"
              autoFocus
              value={formData.name}
              onChange={handleChange}
              disabled={mutation.isPending}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              id="username"
              label="Username"
              name="username"
              autoComplete="username"
              value={formData.username}
              onChange={handleChange}
              disabled={mutation.isPending}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="new-password"
              value={formData.password}
              onChange={handleChange}
              disabled={mutation.isPending}
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2, py: 1.5 }}
              disabled={mutation.isPending}
            >
              {mutation.isPending ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                'Sign Up'
              )}
            </Button>
          </Box>
          {mutation.isError && (
            <Typography color="red">
              {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
              {(mutation.error as any).response?.data?.message}
            </Typography>
          )}
        </Paper>
      </Box>
    </Container>
  );
}
