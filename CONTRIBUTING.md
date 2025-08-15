# Contributing to ZeePass

🎉 Thank you for considering contributing to ZeePass! We welcome contributions from developers who share our commitment to privacy, security, and open-source software.

## 🌟 Ways to Contribute

### 🐛 Bug Reports
- Search existing issues first to avoid duplicates
- Use the bug report template when creating new issues
- Include detailed steps to reproduce the issue
- Provide system information (OS, Go version, browser)
- Add screenshots or logs when helpful

### ✨ Feature Requests
- Check if the feature has already been requested
- Clearly describe the problem you're trying to solve
- Explain how this feature would benefit ZeePass users
- Consider security implications for any new features

### 🔧 Code Contributions
- Bug fixes and security improvements
- Performance optimizations
- New encryption features or tools
- UI/UX improvements
- Documentation improvements
- Test coverage enhancements

### 📖 Documentation
- Fix typos and grammar errors
- Improve code examples
- Add missing documentation
- Translate documentation (future)

## 🚀 Getting Started

### Prerequisites
- Go 1.24.2 or higher
- Redis server
- Git
- Basic understanding of cryptography (for security-related contributions)

### Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**
   ```bash
   git clone https://github.com/your-username/zeepass.git
   cd zeepass/src
   ```

3. **Set up upstream remote**
   ```bash
   git remote add upstream https://github.com/anazri/zeepass.git
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Start Redis server**
   ```bash
   # macOS with Homebrew
   brew services start redis
   
   # Ubuntu/Debian
   sudo systemctl start redis-server
   ```

6. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

7. **Verify setup**
   - Open http://localhost:8080
   - Test basic functionality

## 📝 Contribution Process

### 1. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. Make Your Changes
- Write clean, readable code
- Follow existing code style and conventions
- Add comments for complex logic
- Include tests for new functionality

### 3. Test Your Changes
```bash
# Run tests (when available)
go test ./...

# Test manually with the application
go run cmd/server/main.go
```

### 4. Commit Your Changes
Use clear, descriptive commit messages:
```bash
git add .
git commit -m "Add AES-256 key rotation feature

- Implement automatic key rotation every 30 days
- Add configuration option for rotation interval
- Include backward compatibility for existing encrypted data
- Add comprehensive tests for key rotation logic"
```

### 5. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
```

Then create a pull request through GitHub's interface.

## 📋 Pull Request Guidelines

### Before Submitting
- [ ] Code follows the project's style conventions
- [ ] Tests pass (when available)
- [ ] Documentation is updated if needed
- [ ] Security implications have been considered
- [ ] No sensitive information is exposed in logs or code

### PR Description Template
```markdown
## Summary
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] Security improvement

## Testing
- [ ] Tested locally
- [ ] Added/updated tests
- [ ] Manual testing performed

## Security Considerations
- [ ] No sensitive data exposed
- [ ] Cryptographic changes reviewed
- [ ] Input validation implemented

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

## 🔒 Security Considerations

**⚠️ CRITICAL: Security-First Approach**

Since ZeePass handles sensitive data and cryptographic operations:

- **Never commit secrets** or hardcoded keys
- **Review cryptographic implementations** thoroughly
- **Follow security best practices** for input validation
- **Consider timing attacks** and side-channel vulnerabilities
- **Test with malicious inputs** to prevent exploitation
- **Report security vulnerabilities** privately to security@moonkite.io

### Security Review Process
1. All cryptography-related PRs require security review
2. Performance changes affecting crypto operations need review
3. Input validation changes require thorough testing
4. UI changes handling sensitive data need security assessment

## 🎯 Code Style Guidelines

### Go Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Keep functions focused and small
- Add comments for public functions
- Handle errors appropriately
- Use context for cancellation and timeouts

### Frontend Code Style
- Use semantic HTML elements
- Follow TailwindCSS conventions
- Keep JavaScript minimal and focused
- Use HTMX attributes consistently
- Ensure accessibility compliance

### Example Code Structure
```go
// ✅ Good
func encryptData(ctx context.Context, data []byte, key []byte) ([]byte, error) {
    if len(data) == 0 {
        return nil, errors.New("data cannot be empty")
    }
    
    // Implementation with proper error handling
    encrypted, err := performEncryption(data, key)
    if err != nil {
        return nil, fmt.Errorf("encryption failed: %w", err)
    }
    
    return encrypted, nil
}

// ❌ Avoid
func encrypt(d []byte, k []byte) []byte {
    // Missing validation, error handling, and context
    return performEncryption(d, k)
}
```

## 🧪 Testing Guidelines

### Test Categories
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Security Tests**: Test for vulnerabilities and edge cases
4. **Performance Tests**: Ensure acceptable performance

### Writing Good Tests
```go
func TestEncryptDecryptRoundTrip(t *testing.T) {
    testCases := []struct {
        name string
        data []byte
        key  []byte
    }{
        {"empty data", []byte(""), generateKey()},
        {"small data", []byte("hello"), generateKey()},
        {"large data", make([]byte, 1024*1024), generateKey()},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            encrypted, err := encryptData(context.Background(), tc.data, tc.key)
            require.NoError(t, err)
            
            decrypted, err := decryptData(context.Background(), encrypted, tc.key)
            require.NoError(t, err)
            require.Equal(t, tc.data, decrypted)
        })
    }
}
```

## 🏷️ Issue Labels

We use the following labels to categorize issues:

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `security` - Security-related issues
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `performance` - Performance improvements
- `ui/ux` - User interface and experience
- `crypto` - Cryptography-related changes

## 🎖️ Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Project documentation (with permission)
- Annual contributor appreciation

## 📞 Getting Help

- **General questions**: Open a GitHub discussion
- **Bug reports**: Create an issue with the bug template
- **Security concerns**: Email security@moonkite.io
- **Feature discussions**: Open an issue or discussion
- **Development setup**: Check existing issues or create new one

## 📜 Code of Conduct

### Our Standards
- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Respect different perspectives
- Prioritize security and privacy

### Unacceptable Behavior
- Harassment or discrimination
- Sharing others' private information
- Spam or off-topic content
- Disrespectful or unprofessional conduct

## 🚀 Development Roadmap

Check our [GitHub Projects](https://github.com/anazri/zeepass/projects) for:
- Current development priorities
- Upcoming features
- Community-requested enhancements
- Good first issues for new contributors

## 📄 License

By contributing to ZeePass, you agree that your contributions will be licensed under the MIT License.

---

**🔐 Remember: ZeePass is a security tool. Every contribution should maintain the highest standards of security and privacy.**

Thank you for helping make ZeePass better for everyone! 🙏

---

*For commercial support, deployment services, or consulting: contact@moonkite.io*