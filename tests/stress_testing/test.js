import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend, Counter } from 'k6/metrics';
import { randomIntBetween, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const trendDeactivate = new Trend('duration_deactivate');
const trendMerge = new Trend('duration_merge');
const trendCreatePR = new Trend('duration_create_pr');
const trendStats = new Trend('duration_stats');

const countDeactivate = new Counter('reqs_deactivate');
const countMerge = new Counter('reqs_merge');
const countCreatePR = new Counter('reqs_create_pr');

export const options = {
  stages: [
    { duration: '5s', target: 5 },   
    { duration: '1m', target: 5 },   
    { duration: '5s', target: 0 },   
  ],
  
  thresholds: {
    http_req_failed: ['rate<0.001'], 
    
    'duration_deactivate': ['p(95)<300'], 
    'duration_merge': ['p(95)<300'],
    'duration_create_pr': ['p(95)<300'],
    'duration_stats': ['p(95)<300'],
  },
};

const BASE_URL = 'http://localhost:8080';
const TEAMS_COUNT = 20;
const USERS_PER_TEAM = 10; 

export function setup() {
  console.log(`Creating ${TEAMS_COUNT} teams with ${USERS_PER_TEAM} users each...`);
  
  let allUserIds = [];
  let allTeamNames = [];

  for (let t = 1; t <= TEAMS_COUNT; t++) {
    const teamName = `Team_${t}`;
    const members = [];
    
    for (let u = 1; u <= USERS_PER_TEAM; u++) {
      const userId = `user_t${t}_u${u}`;
      members.push({
        user_id: userId,
        username: `User ${t}-${u}`,
        is_active: true,
      });
      allUserIds.push(userId);
    }

    const res = http.post(`${BASE_URL}/team/add`, JSON.stringify({
      team_name: teamName,
      members: members,
    }), { headers: { 'Content-Type': 'application/json' } });
    
    const isOk = check(res, { 'setup team created': (r) => r.status === 200 || r.status === 201 });
    if (!isOk) console.error(`Setup Failed for ${teamName}: ${res.status}`);
    
    allTeamNames.push(teamName);
  }
  console.log('Setup complete.');
  return { allUserIds, allTeamNames };
}

export default function (data) {
  const rand = Math.random();

  if (rand < 0.4) {
    group('Create PR', () => {
      const authorId = randomItem(data.allUserIds);
      const prId = `pr_${__VU}_${Date.now()}_${Math.random().toString(36).substring(7)}`; 
      
      const res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify({
        pull_request_id: prId,
        pull_request_name: `Feature ${prId}`,
        author_id: authorId,
      }), { headers: { 'Content-Type': 'application/json' } });

      trendCreatePR.add(res.timings.duration);
      countCreatePR.add(1);

      check(res, {
        'pr created': (r) => r.status === 200 || r.status === 201,
      });
    });
  } 

  else if (rand < 0.7) {
    group('Merge', () => {
      const authorId = randomItem(data.allUserIds);
      const prId = `merge_${__VU}_${Date.now()}_${Math.random().toString(36).substring(7)}`;

      http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify({
        pull_request_id: prId, pull_request_name: "M", author_id: authorId
      }), { headers: { 'Content-Type': 'application/json' } });

      const mergeRes = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify({
          pull_request_id: prId
      }), { headers: { 'Content-Type': 'application/json' } });

      trendMerge.add(mergeRes.timings.duration);
      countMerge.add(1);

      check(mergeRes, { 'merge success': (r) => r.status === 200 || r.status === 201 });
    });
  }

  else if (rand < 0.8) {
    group('Deactivate', () => {
      const u1 = randomItem(data.allUserIds);
      const u2 = randomItem(data.allUserIds);

      const res = http.post(`${BASE_URL}/users/deactivate`, JSON.stringify({
        user_ids: [u1, u2]
      }), { headers: { 'Content-Type': 'application/json' } });
      trendDeactivate.add(res.timings.duration); 
      countDeactivate.add(1);                   

      check(res, { 'deactivate success': (r) => r.status === 200 });
    });
  }

  else {
    const res = http.get(`${BASE_URL}/stats`);
    trendStats.add(res.timings.duration);
    check(res, { 'stats 200': (r) => r.status === 200 });
  }

  sleep(1); 
}
