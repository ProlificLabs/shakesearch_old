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
    await page.goto('http://localhost:3002');
  });

  after(async () => {
    await browser.close();
  });

  it('should return no results for "Luke, I am your father"', async () => {
    const query = 'Luke, I am your father';
    await page.evaluate(() => document.getElementById('query').value = ''); // Reset the input field
    await page.type('#query', query);
    await page.click('button[type="submit"]');
    await page.waitForTimeout(1000);
  
    const results = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
        return Array.from(rows, row => row.textContent.trim());

    });
  
    assert.isEmpty(results, 'Search results should be empty');
  });

  it('should return search results for "romeo, wherefore art thou"', async () => {
    const query = 'romeo, wherefore art thou';
    await page.evaluate(() => document.getElementById('query').value = ''); // Reset the input field
    await page.type('#query', query);
    await page.click('button[type="submit"]');
    await page.waitForSelector('#table-body tr');

    const results = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
        return Array.from(rows, row => row.textContent.trim());

    });

    assert.isNotEmpty(results, 'Search results should not be empty');
    assert.include(results.join(' ').toLowerCase(), query.toLowerCase(), 'Search results should contain the query');
  });

  it('should load more results for "horse" when clicking "Load More"', async () => {
    const query = 'horse';
    await page.evaluate(() => document.getElementById('query').value = ''); // Reset the input field
    await page.type('#query', query);
    await page.click('button[type="submit"]');
    await page.waitForSelector('#table-body tr');

    const initialResults = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
        return Array.from(rows, row => row.textContent.trim());

    });

    await page.click('#load-more');
    await page.waitForTimeout(1000);

    const updatedResults = await page.evaluate(() => {
      const rows = document.querySelectorAll('#table-body tr');
        return Array.from(rows, row => row.textContent.trim());

    });

    assert.isAbove(updatedResults.length, initialResults.length, 'More results should be added to the results table');
  });

});