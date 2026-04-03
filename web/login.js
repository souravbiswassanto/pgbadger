// login.js
// Read cookie helper
function readCookie(name) {
  const v = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')
  return v ? v.pop() : ''
}

document.getElementById('loginForm').addEventListener('submit', async (e) => {
  e.preventDefault()
  const username = document.getElementById('username').value
  const password = document.getElementById('password').value

  // Double-submit CSRF: send X-CSRF-Token header with value from csrf_token cookie
  const csrf = readCookie('csrf_token')

  const resp = await fetch('/api/v1/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrf
    },
    body: JSON.stringify({ username, password }),
    credentials: 'same-origin'
  })

  if (resp.ok) {
    alert('login successful')
    window.location.href = '/'
  } else {
    const body = await resp.json().catch(()=>({}))
    alert('login failed: ' + (body.error || resp.status))
  }
})
// external login script (keeps JS out of HTML to enable strict CSP)
function readCookie(name) {
  const v = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')
  return v ? v.pop() : ''
}

document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('loginForm')
  if (!form) return

  form.addEventListener('submit', async (e) => {
    e.preventDefault()
    const username = document.getElementById('username').value
    const password = document.getElementById('password').value

    // Double-submit CSRF: send X-CSRF-Token header with value from csrf_token cookie
    const csrf = readCookie('csrf_token')

    try {
      const resp = await fetch('/api/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrf
        },
        body: JSON.stringify({ username, password }),
        credentials: 'same-origin'
      })

      if (resp.ok) {
        alert('login successful')
        window.location.href = '/'
      } else {
        const body = await resp.json().catch(()=>({}))
        alert('login failed: ' + (body.error || resp.status))
      }
    } catch (err) {
      alert('network or server error')
    }
  })
})
