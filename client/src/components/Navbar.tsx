import {
  AppBar,
  Box,
  Toolbar,
  Typography,
  Button,
  Container,
} from '@mui/material';
import AdbIcon from '@mui/icons-material/Adb';
import { Link, useNavigate } from 'react-router-dom';

const Navbar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <AppBar position="static">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Box display="flex" alignItems="center" sx={{ marginRight: 4 }}>
            <AdbIcon sx={{ mr: 1 }} />
            <Typography
              variant="h6"
              noWrap
              component={Link}
              to="/dashboard"
              sx={{
                fontFamily: 'monospace',
                fontWeight: 700,
                letterSpacing: '.3rem',
                color: 'inherit',
                textDecoration: 'none',
              }}
            >
              Workflow Automation System
            </Typography>
          </Box>

          <Box sx={{ flexGrow: 1 }}>
            <Button
              component={Link}
              to="/community"
              sx={{ my: 2, color: 'white', display: 'block' }}
            >
              Community
            </Button>
          </Box>

          <Box sx={{ flexGrow: 1 }}>
            <Button
              component={Link}
              to="/connections"
              sx={{ my: 2, color: 'white', display: 'block' }}
            >
              Connections
            </Button>
          </Box>

          <Box sx={{ flexGrow: 1 }}>
            <Button
              component={Link}
              to="/templates"
              sx={{ my: 2, color: 'white', display: 'block' }}
            >
              Templates
            </Button>
          </Box>

          <Box sx={{ flexGrow: 0 }}>
            <Button color="inherit" onClick={handleLogout}>
              Logout
            </Button>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Navbar;
