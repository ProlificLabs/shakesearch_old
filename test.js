const assert = require('chai').assert;
const puppeteer = require('puppeteer-core');

describe('ShakeSearch', () => {
  let browser;
  let page;

  before(async () => {
    browser = await puppeteer.launch({
      executablePath: '/usr/bin/chromium-browser',
      args: ['--no-sandbox'],
    });
    page = await browser.newPage();
    await page.goto('http://localhost:3001');
  });

  after(async () => {
    await browser.close();
  });

  it('should return no results for "Luke, I am your father"', async () => {
    const query = 'Luke, I am your father';
    await page.type('#query', query);
    await page.click('button[type="submit"]');
    await page.waitForTimeout(1000);
  
    const results = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
      return Array.from(rows, row => row.innerText);
    });
  
    assert.isEmpty(results, 'Search results should be empty');
  });

  it('should return search results for "wherefore art thou romeo"', async () => {
    const query = 'wherefore art thou romeo';
    await page.type('#query', query);
    await page.click('button[type="submit"]');
    await page.waitForSelector('#table-body tr');

    const results = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
      return Array.from(rows, row => row.innerText);
    });

    assert.isNotEmpty(results, 'Search results should not be empty');
    assert.include(results.join(' ').toLowerCase(), query.toLowerCase(), 'Search results should contain the query');
  });

  it('should submit search by pressing "Enter"', async () => {
    const query = 'Hamlet';
    await page.type('#query', query);
    await page.keyboard.press('Enter');
    await page.waitForSelector('#table-body tr');

    const results = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
      return Array.from(rows, row => row.innerText);
    });

    assert.isNotEmpty(results, 'Search results should not be empty');
  });
});