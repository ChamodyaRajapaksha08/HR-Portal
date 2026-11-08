<template>
  <div class="login-container">
    <div class="login-window">
      
      <div class="brand-panel">
        <div class="brand-background-image" style="background-image: url('/login_bg.jpg');">
          <div class="green-transparent-overlay"></div>
        </div>

        <div class="brand-contents-wrapper">
          <div class="brand-welcome-center">
            <h1>Welcome Back!</h1>
            <p>Sign in to continue managing your organization efficiently.</p>
          </div>
          <p class="brand-tagline">Streamlining workforce management with intelligence.</p>
        </div>
      </div>

      <div class="form-panel">
        <div class="form-center-box">
          
          <div class="logo-wrapper">
            <img src="/system_logo.png" alt="Logo" class="login-page-logo" />
          </div>
          
          <h2>HR Management System</h2>
          <p class="form-subtitle">Sign in to your account</p>

          <form @submit.prevent="handleLogin">
            <div class="form-group">
              <label>Email Address</label>
              <div class="input-with-icon">
                <span class="input-icon">✉️</span>
                <input type="email" v-model="email" placeholder="Enter your email" class="input-field" required />
              </div>
            </div>

            <div class="form-group">
              <label>Password</label>
              <div class="input-with-icon">
                <span class="input-icon">🔒</span>
                <input type="password" v-model="password" placeholder="Enter your password" class="input-field" required />
              </div>
            </div>

            <div class="form-options">
              <label class="remember-me">
                <input v-model="rememberMe" type="checkbox" />
                <span>Remember me</span>
              </label>
              <a class="forgot-link" :href="`mailto:admin@itum.com?subject=Password%20reset%20request&body=Please%20help%20me%20reset%20the%20password%20for%20${encodeURIComponent(email)}`">Forgot Password?</a>
            </div>

            <p v-if="errorMessage" style="color: #ef4444; font-size: 12px; text-align: center; margin-bottom: 10px; font-weight: 600;">
              {{ errorMessage }}
            </p>

            <button type="submit" class="login-submit-btn">
              Sign In →
            </button>
          </form>

          <p class="contact-admin-text">
            Don't have an account? <a class="contact-link" href="mailto:admin@itum.com?subject=HRM%20account%20request">Contact Admin</a>
          </p>
        </div>
        <p class="copyright-text">Copyright © HR Management System. All rights reserved.</p>
      </div>

    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { apiPost } from '../api';

const router = useRouter();
const route = useRoute();

const email = ref('');
const password = ref('');
const errorMessage = ref('');
const rememberMe = ref(false);

onMounted(() => {
  if (route.query.sessionExpired === '1') {
    errorMessage.value = 'Session expired, please log in again';
  }
  const savedEmail = localStorage.getItem('rememberedEmail');
  if (savedEmail) {
    email.value = savedEmail;
    rememberMe.value = true;
  }
});

const handleLogin = async () => {
  errorMessage.value = '';

  try {
    const data = await apiPost('/api/login', {
      email: email.value,
      password: password.value
    });

    if (!data.token || !data.user) throw new Error('The server returned an invalid login response.');
    localStorage.setItem('authToken', data.token);
    const user = { ...data.user, full_name: data.user.name };
    localStorage.setItem('authUser', JSON.stringify(user));
    if (rememberMe.value) {
      localStorage.setItem('rememberedEmail', email.value);
    } else {
      localStorage.removeItem('rememberedEmail');
    }
    const redirect = router.currentRoute.value.query.redirect;
    const home = '/dashboard';
    router.replace(
      typeof redirect === 'string' && redirect.startsWith('/') && !redirect.startsWith('//')
        ? redirect
        : home
    );
  } catch (error) {
    errorMessage.value = error.message;
    console.error('Login failed.');
  }
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght=400;500;600;700&display=swap');

.login-container { 
  display: flex !important; 
  justify-content: center !important; 
  align-items: center !important; 
  min-height: 100vh !important; 
  background-color: #f1f5f9 !important; 
  font-family: 'Inter', sans-serif !important; 
  padding: 20px !important; 
  box-sizing: border-box !important; 
}

.login-window { 
  display: grid !important; 
  grid-template-columns: 1.4fr 1fr !important; 
  width: 100% !important; 
  max-width: 1060px !important; 
  height: 540px !important; 
  background: white !important; 
  border-radius: 20px !important; 
  overflow: hidden !important; 
  box-shadow: 0 25px 50px -12px rgba(6, 44, 34, 0.1) !important; 
  border: 1px solid #e2e8f0 !important; 
}

.brand-panel { position: relative !important; overflow: hidden !important; height: 100% !important; }
.brand-background-image { position: absolute !important; top: 0 !important; left: 0 !important; right: 0 !important; bottom: 0 !important; background-size: cover !important; background-position: center !important; }
.green-transparent-overlay { position: absolute !important; top: 0 !important; left: 0 !important; right: 0 !important; bottom: 0 !important; background: linear-gradient(to bottom, rgba(6, 36, 27, 0.82) 0%, rgba(18, 61, 46, 0.62) 50%, rgba(6, 36, 27, 0.92) 100%) !important; }
.brand-contents-wrapper { position: relative !important; z-index: 10 !important; height: 100% !important; padding: 40px 40px 28px !important; display: flex !important; flex-direction: column !important; justify-content: space-between !important; }
.brand-welcome-center { text-align: left !important; margin-top: 10px !important; }
.login-container .brand-panel .brand-welcome-center h1 { font-family: 'Inter', sans-serif !important; font-size: 40px !important; font-weight: 700 !important; color: #ffffff !important; text-shadow: 0 2px 8px rgba(2, 6, 23, 0.65) !important; margin: 0 0 12px 0 !important; }
.login-container .brand-panel .brand-welcome-center p { font-family: 'Inter', sans-serif !important; font-size: 13px !important; color: #e2e8f0 !important; max-width: 280px !important; margin: 0 !important; line-height: 1.6 !important; }
.login-container .brand-panel .brand-tagline { font-family: 'Inter', sans-serif !important; margin: auto 0 0 !important; font-size: 12px !important; font-weight: 500 !important; line-height: 1.5 !important; color: rgba(255, 255, 255, 0.8) !important; text-shadow: 0 1px 5px rgba(2, 6, 23, 0.5) !important; }

.form-panel { padding: 40px 25px !important; display: flex !important; flex-direction: column !important; justify-content: space-between !important; align-items: center !important; background-color: #f0faf5 !important; height: 100% !important; }
.form-center-box { width: 100% !important; max-width: 320px !important; margin: auto 0 !important; }
.login-page-logo { display: block !important; margin: -15px auto 15px auto !important; width: 65px !important; height: auto !important; }
.form-center-box h2 { font-family: 'Inter', sans-serif !important; font-size: 20px !important; margin: 0 0 4px 0 !important; text-align: center !important; color: #18352b !important; font-weight: 700 !important; letter-spacing: -0.5px !important; }
.form-subtitle { font-family: 'Inter', sans-serif !important; text-align: center !important; font-size: 12px !important; color: #475569 !important; margin: 0 0 15px 0 !important; }

.form-group { margin-bottom: 12px !important; }
label { font-family: 'Inter', sans-serif !important; display: block !important; font-size: 11px !important; font-weight: 600 !important; color: #18352b !important; margin-bottom: 4px !important; text-transform: uppercase !important; letter-spacing: 0.5px !important; }
.input-with-icon { position: relative !important; display: flex !important; align-items: center !important; }
.input-icon { position: absolute !important; left: 14px !important; color: #64748b !important; font-size: 12px !important; }
.input-field { width: 100% !important; padding: 10px 14px 10px 38px !important; border: 1px solid #d7e5df !important; border-radius: 8px !important; font-size: 13px !important; outline: none !important; background-color: #ffffff !important; color: #18352b !important; font-family: 'Inter', sans-serif !important; }

.form-options { display: flex !important; justify-content: space-between !important; align-items: center !important; font-size: 12px !important; margin-bottom: 16px !important; }
.remember-me { display: flex !important; align-items: center !important; gap: 6px !important; color: #334155 !important; cursor: pointer !important; font-weight: 500 !important; }
.forgot-link { color: #047857 !important; font-weight: 600 !important; cursor: pointer !important; text-decoration: none !important; }

.login-submit-btn { width: 100% !important; padding: 12px !important; background-color: #047857 !important; color: white !important; border: none !important; border-radius: 8px !important; font-size: 14px !important; font-weight: 600 !important; cursor: pointer !important; transition: 0.2s !important; font-family: 'Inter', sans-serif !important; text-align: center !important; display: block !important; box-sizing: border-box !important; }
.login-submit-btn:hover { background-color: #059669 !important; }

.contact-admin-text { font-family: 'Inter', sans-serif !important; text-align: center !important; font-size: 12px !important; color: #475569 !important; margin: 16px 0 0 !important; }
.contact-link { color: #047857 !important; font-weight: 600 !important; cursor: pointer !important; }
.copyright-text { font-family: 'Inter', sans-serif !important; font-size: 10px !important; color: #64748b !important; text-align: center !important; margin: 0 !important; }
</style>