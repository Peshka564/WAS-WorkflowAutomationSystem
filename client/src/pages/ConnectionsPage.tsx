import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  Box,
  Button,
  Typography,
  Paper,
  Chip,
  Container,
  Snackbar,
  Alert,
} from '@mui/material';
import GoogleIcon from '@mui/icons-material/Google';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import axios from 'axios';

interface Connection {
  service: string;
  connected: boolean;
}

export function ConnectionsPage() {
  const [connections, setConnections] = useState<Connection[]>([]);
  const [searchParams, setSearchParams] = useSearchParams();
  //   const navigate = useNavigate();

  // Toast State
  const [toast, setToast] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error';
  }>({
    open: false,
    message: '',
    severity: 'success',
  });

  useEffect(() => {
    fetchConnections();
    handleRedirectParams();
  }, []);

  const handleRedirectParams = () => {
    const status = searchParams.get('status');

    if (status === 'success') {
      setToast({
        open: true,
        message: 'Successfully connected to Gmail!',
        severity: 'success',
      });
      // Clean the URL so refresh doesn't trigger it again
      setSearchParams({});
    } else if (status === 'error') {
      setToast({
        open: true,
        message: 'Failed to connect. Please try again.',
        severity: 'error',
      });
      setSearchParams({});
    }
  };

  const fetchConnections = async () => {
    const token = localStorage.getItem('token');
    try {
      const res = await axios.get('http://localhost:3000/api/connections', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setConnections(res.data);
    } catch (e) {
      console.error(e);
    }
  };

  const handleConnect = async () => {
    try {
      const token = localStorage.getItem('token');

      const res = await axios.get(
        'http://localhost:3000/api/auth/google/login',
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (res.data.url) {
        window.location.href = res.data.url;
      }
    } catch (error) {
      console.error('Failed to start OAuth flow', error);
    }
  };

  const gmailConnection = connections.find(c => c.service === 'gmail');

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      <Typography variant="h4" gutterBottom fontWeight="bold">
        Integrations & Connections
      </Typography>
      <Typography color="textSecondary" sx={{ mb: 4 }}>
        Manage your third-party connections globally.
      </Typography>

      <Paper elevation={2} sx={{ p: 3 }}>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Box display="flex" alignItems="center" gap={2}>
            <GoogleIcon sx={{ color: '#DB4437', fontSize: 40 }} />
            <Box>
              <Typography variant="h6">Gmail</Typography>
              <Typography variant="body2" color="textSecondary">
                Allows sending emails via your Google account.
              </Typography>
            </Box>
          </Box>

          <Box>
            {gmailConnection ? (
              <Box display="flex" alignItems="center" gap={1}>
                <Chip
                  icon={<CheckCircleIcon />}
                  label="Connected"
                  color="success"
                  variant="outlined"
                />
                {/* Optional: Add Disconnect Logic later */}
                {/* <Button color="error" size="small">Disconnect</Button> */}
              </Box>
            ) : (
              <Button
                variant="contained"
                onClick={handleConnect}
                sx={{ bgcolor: '#DB4437', '&:hover': { bgcolor: '#C53929' } }}
              >
                Connect with Google
              </Button>
            )}
          </Box>
        </Box>
      </Paper>

      {/* Success/Error Toast */}
      <Snackbar
        open={toast.open}
        autoHideDuration={6000}
        onClose={() => setToast(prev => ({ ...prev, open: false }))}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert
          severity={toast.severity}
          variant="filled"
          onClose={() => setToast(prev => ({ ...prev, open: false }))}
        >
          {toast.message}
        </Alert>
      </Snackbar>
    </Container>
  );
}
